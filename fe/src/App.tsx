import React, { useState, useEffect, useMemo, useCallback, FunctionComponent } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import { AgGridReact } from 'ag-grid-react';
import { ColumnApi, GridReadyEvent, ModuleRegistry, CellEditingStartedEvent, CellEditingStoppedEvent, CellValueChangedEvent, ColDef } from 'ag-grid-community';

import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';

import {
    GET_CELLS_BY_SPREADSHEET_ID,
    GET_SPREADSHEET,
    UPDATE_CELL_BY_SPREADSHEET_ID_COLUMN_AND_ROW,
    UPDATE_SPREADSHEET
} from './graphqlQueries';
import {Cell} from "./__generated__/graphql";
import {ClientSideRowModelModule} from "@ag-grid-community/client-side-row-model";
import './App.css';
import {generateColumnFields, convertCellsToRowData} from "./utils";


interface RowData {
    rowIndex: string;
    rawValue: string;
    computedValue: string;
    editMode: number;
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

    const [updateSpreadsheet, { data: dataUpdateSpreadsheet, loading: loadingUpdateSpreadsheet, error: errorUpdateSpreadsheet }] = useMutation(
        UPDATE_SPREADSHEET,
        {
            refetchQueries: [{query: GET_SPREADSHEET}, {query: GET_CELLS_BY_SPREADSHEET_ID}],
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
        if (data && dataSpreadsheet) {
            const preparedRowData = convertCellsToRowData(data.getCellsBySpreadsheetId, dataSpreadsheet.getSpreadsheet.rowCount);
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
        if (data && dataSpreadsheet) {
            const preparedRowData = convertCellsToRowData(data.getCellsBySpreadsheetId, dataSpreadsheet.getSpreadsheet.rowCount);
            setRowData(preparedRowData);
        } else {
            console.log('data is null');
        }
    }, [data, dataSpreadsheet]);


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

        // if they have edited a cell in a new row or column we have to update the sheet first
        // the UX here is cludgy, but it's a demo
        if (dataSpreadsheet.getSpreadsheet.rowCount < (rowIndex||0) + 1 || dataSpreadsheet.getSpreadsheet.columnCount < column.getInstanceId()) {
            const newRowCount = Math.max((rowIndex||0) + 1, dataSpreadsheet.getSpreadsheet.rowCount);
            const newColumnCount = Math.max(column.getInstanceId(), dataSpreadsheet.getSpreadsheet.columnCount);
            updateSpreadsheet({variables: {
                    id: dataSpreadsheet.getSpreadsheet.id,
                    rowCount: newRowCount,
                    columnCount: newColumnCount
                }})

            if (newColumnCount > dataSpreadsheet.getSpreadsheet.columnCount) {
                const columnDefs = generateColumnFields(newColumnCount);
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
            }
        }
        updateCell({
            variables: {
                spreadsheetId: dataSpreadsheet.getSpreadsheet.id,
                columnIndex: event.column.getInstanceId() - 1,
                rowIndex: event.node.rowIndex || 0,
                rawValue:  event.newValue,
            },
        });
        const columnCode = column.getColId();
        console.log("EventCellEditingStopped,", event);
        if (columnCode && event.data && event.data[columnCode]) event.data[columnCode].editMode = 1;
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
                    ></AgGridReact>
                </div>
            </div>
        </div>
    );
};

export default App;
