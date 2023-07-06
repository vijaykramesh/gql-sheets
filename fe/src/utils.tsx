import { ColDef } from 'ag-grid-community';
import { Cell } from './__generated__/graphql';
import ReactMarkdown from 'react-markdown';
import React from 'react';
import { renderToStaticMarkup } from 'react-dom/server';

export function columnCodeFromColumnIndex(columnIndex: number): string {
    if (columnIndex < 0) {
        return '';
    }

    let column = '';
    while (columnIndex >= 0) {
        const remainder = columnIndex % 26;
        column = String.fromCharCode(remainder + 65) + column;
        // eslint-disable-next-line no-param-reassign
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

                const cellData = params.data[columnCode];
                const rawCellValue = params.value;
                if (cellData) {
                    if (cellData.editMode && cellData.editMode === 0) {
                        return cellData.rawValue;
                    }
                    const markdownContent = renderToStaticMarkup(<ReactMarkdown>{cellData.computedValue}</ReactMarkdown>);
                    return `<span dangerouslySetInnerHTML={{ __html: ${markdownContent} }} />`;

                    // return cellData.computedValue;
                }
                return rawCellValue;
            },
            valueParser: (params) => {
                return params.newValue;
            },
            valueGetter: (params) => {
                if (params.data[columnCode]) {
                    const markdownContent = renderToStaticMarkup(<ReactMarkdown>{params.data[columnCode].computedValue}</ReactMarkdown>);

                    return params.data[columnCode].editMode !== 0 ?
                        `<span dangerouslySetInnerHTML={{ __html: ${markdownContent} }} />` :
                        params.data[columnCode].rawValue;
                }
                return '';
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

export function convertCellsToRowData(data: Cell[], rowCount: number, editModeColumnAndRow: {
    column: string;
    row: number
} | null): { [key: string]: any }[] {
    const maxRowIndex = Math.max(rowCount - 1, Math.max(...data.map((cell) => {
        return cell.rowIndex;
    })));
    const maxColumnIndex = Math.max(...data.map((cell) => {
        return cell.columnIndex;
    }));

    const result: { [key: string]: object }[] = [];
    for (let i = 0; i <= maxRowIndex; i++) {
        const row: { [key: string]: any } = {};

        for (let j = 0; j <= maxColumnIndex; j++) {
            const cell = data.find((c) => {
                return c.rowIndex === i && c.columnIndex === j;
            });
            const columnCode = columnCodeFromColumnIndex(j);
            let editMode = 0;
            if (editModeColumnAndRow && editModeColumnAndRow.row === i && editModeColumnAndRow.column === columnCode) {
                editMode = 1;
            }
            row[columnCode] = cell ?
                {
                    rowIndex: (i + 1).toString(),
                    rawValue: cell.rawValue.toString(),
                    computedValue: (cell.computedValue || '').toString(),
                    editMode: editMode,
                } :
                {
                    rowIndex: (i + 1).toString(),
                    rawValue: '',
                    computedValue: '',
                    editMode: 0,
                };
        }
        row.rowIndex = (i + 1).toString();
        result.push(row);
    }
    const newRow: { [key: string]: any } = {};
    newRow.rowIndex = (maxRowIndex + 2).toString();
    result.push(newRow);
    return result;
}
