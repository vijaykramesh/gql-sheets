import {ColDef} from "ag-grid-community";
import {Cell} from "./__generated__/graphql";

export function columnCodeFromColumnIndex(columnIndex: number): string {
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

export function generateColumnFields(maxColumnIndex: number): ColDef[] {
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

export
function convertCellsToRowData(data: Cell[], rowCount: number): { [key: string]: any }[] {
    const maxRowIndex = Math.max(rowCount-1,Math.max(...data.map((cell) => cell.rowIndex)))
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
