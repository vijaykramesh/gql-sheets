import React, { useState, useEffect, useMemo, useCallback, FunctionComponent } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import { AgGridReact } from 'ag-grid-react';
import { ColumnApi, GridReadyEvent, ModuleRegistry, CellEditingStartedEvent, CellEditingStoppedEvent, CellValueChangedEvent, ColDef } from 'ag-grid-community';

import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';

import { GET_CELLS_BY_SPREADSHEET_ID, GET_SPREADSHEET, UPDATE_CELL_BY_SPREADSHEET_ID_COLUMN_AND_ROW } from './graphqlQueries';
import {Cell} from "./__generated__/graphql";
import {ClientSideRowModelModule} from "@ag-grid-community/client-side-row-model";

function columnCodeFromColumnIndex(columnIndex: number): string {
    if (columnIndex < 0) {
        return '';
    }

    let column = '';
    while (columnIndex >= 0) {
        const remainder = columnIndex % 26;
        column = String.fromCharCode(remainder + 65) + column;
        columnIndex = Math.floor(columnIndex / 26) - 1;

        if (columnIndex < 0) {
            break;
        }
    }

    return column;
}

function generateColumnFields(maxColumnIndex: number): ColDef[] {
    const fields: ColDef[] = [];

    for (let i = 0; i <= maxColumnIndex; i++) {
        const columnCode = columnCodeFromColumnIndex(i);
        const field: ColDef = {
            field: columnCode,
            valueFormatter: (params) => {
                if (params.value === undefined) {
                    return '';
                }

                const cellData = params.value;
                if (cellData.editMode && cellData.editMode === 2) {
                    return cellData.rawValue;
                } else {
                    return cellData.computedValue;
                }
            },
            valueParser: (params) => {
                return params.newValue;
            },
            valueGetter: (params) => {
                if (params.data[columnCode]) {
                    return params.data[columnCode].editMode === 1
                        ? params.data[columnCode].rawValue
                        : params.data[columnCode].computedValue;
                } else {
                    return '';
                }
            },
            valueSetter: (params) => {
                if (params.data[columnCode]) {
                    params.data[columnCode].rawValue = params.newValue;
                }
                return true;
            },
            cellRenderer: 'agAnimateShowChangeCellRenderer',
        };

        fields.push(field);
    }

    return fields;
}

interface RowData {
    rowIndex: string;
    rawValue: string;
    computedValue: string;
    editMode: number;
}

function convertCellsToRowData(data: Cell[]): { [key: string]: any }[] {
    const maxRowIndex = Math.max(...data.map((cell) => cell.rowIndex));
    const maxColumnIndex = Math.max(...data.map((cell) => cell.columnIndex));

    const result: { [key: string]: object }[] = [];
    for (let i = 0; i <= maxRowIndex; i++) {
        const row: { [key: string]: any } = {};

        for (let j = 0; j <= maxColumnIndex; j++) {
            const cell = data.find((c) => c.rowIndex === i && c.columnIndex === j);
            const columnCode = columnCodeFromColumnIndex(j);
            row[columnCode] = cell
                ? {
                    rowIndex: (i + 1).toString(),
                    rawValue: cell.rawValue.toString(),
                    computedValue: (cell.computedValue || '').toString(),
                    editMode: 0,
                }
                : {
                    rowIndex: (i + 1).toString(),
                    rawValue: '',
                    computedValue: '',
                    editMode: 0,
                };
        }
        row['rowIndex'] = (i + 1).toString();
        result.push(row);
    }
    const newRow: { [key: string]: any } = {};
    newRow['rowIndex'] = (maxRowIndex + 2).toString();
    result.push(newRow);
    return result;
}

ModuleRegistry.registerModules([ClientSideRowModelModule]);

const App: FunctionComponent = (): React.ReactElement => {
    const { loading, data } = useQuery(GET_CELLS_BY_SPREADSHEET_ID);
    const { loading: loadingSpreadsheet, data: dataSpreadsheet } = useQuery(GET_SPREADSHEET);
    const [updateCell, { data: dataUpdate, loading: loadingUpdate, error }] = useMutation(
        UPDATE_CELL_BY_SPREADSHEET_ID_COLUMN_AND_ROW,
        {
            refetchQueries: [{ query: GET_CELLS_BY_SPREADSHEET_ID }],
        }
    );

    const containerStyle = useMemo(() => ({ width: '100%', height: '100%' }), []);
    const gridStyle = useMemo(() => ({ height: '500px', width: '90%', padding: '50px' }), []);
    const [rowData, setRowData] = useState<any[]>();
    const [columnDefs, setColumnDefs] = useState<ColDef[]>([
        { headerName: '', field: 'rowIndex', width: 50, sortable: false, resizable: false },
        { field: 'A' },
        { field: 'B' },
        { field: 'C' },
    ]);

    const defaultColDef = useMemo<ColDef>(() => {
        return {
            editable: true,
            cellDataType: false,
        };
    }, [data]);

    const onGridReady = useCallback((params: GridReadyEvent) => {
        if (data) {
            const preparedRowData = convertCellsToRowData(data.getCellsBySpreadsheetId);
            setRowData(preparedRowData);
        } else {
            console.log('data is null');
        }
        if (dataSpreadsheet) {
            const columnDefs = generateColumnFields(dataSpreadsheet.getSpreadsheet.columnCount);
            const firstColumn = [
                {
                    headerName: '',
                    field: 'rowIndex',
                    width: 50,
                    sortable: false,
                    resizable: false,
                    selectable: false,
                    clickable: false,
                    editable: false,
                },
            ];
            columnDefs.unshift(...firstColumn);
            setColumnDefs(columnDefs);
        } else {
            console.log('dataSpreadsheet is null');
        }
    }, [data, dataSpreadsheet]);

    useEffect(() => {
        if (data) {
            const preparedRowData = convertCellsToRowData(data.getCellsBySpreadsheetId);
            setRowData(preparedRowData);
        } else {
            console.log('data is null');
        }
    }, [data]);

    const onCellValueChanged = (event: CellValueChangedEvent) => {
        if (event.colDef.field) {
            updateCell({
                variables: {
                    spreadsheetId: '1',
                    columnIndex: event.column.getInstanceId() - 1,
                    rowIndex: event.node.rowIndex || 0,
                    rawValue: event.data[event.colDef.field] ? event.data[event.colDef.field].rawValue : 'new value',
                },
            });
        }
    };

    const onCellEditingStarted = (event: CellEditingStartedEvent) => {
        const { rowIndex, column } = event;
        const columnCode = column.getColId();
        if (event.data && event.data[columnCode] && event.data[columnCode].editMode === 0) {
            if (event.api) {
                event.api.stopEditing();
            }
        }

        if (event.data && event.data[columnCode]) event.data[columnCode].editMode = 2;
    };

    const onCellEditingStopped = (event: CellEditingStoppedEvent) => {
        const { rowIndex, column } = event;
        const columnCode = column.getColId();
        if (event.data && event.data[columnCode]) event.data[columnCode].editMode = 1;
    };

    const gridOptions = {
        onCellEditingStarted,
        onCellEditingStopped,
        columnHoverHighlight: true,
    };

    if (loading || loadingSpreadsheet) return <p>loading...</p>;
    if (!data && !dataSpreadsheet) return <p>Not found</p>;

    return (
        <div className="App" style={{ height: '100%' }}>
            <header className="App-header" style={{ height: '50px' }}>
                gql-sheets
            </header>
            <div style={containerStyle}>
                <div style={gridStyle} className="ag-theme-alpine">
                    <AgGridReact<Cell>
                        gridOptions={gridOptions}
                        rowData={rowData}
                        columnDefs={columnDefs}
                        defaultColDef={defaultColDef}
                        onGridReady={onGridReady}
                        onCellValueChanged={onCellValueChanged}
                    ></AgGridReact>
                </div>
            </div>
        </div>
    );
};

export default App;
