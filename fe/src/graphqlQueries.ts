import { gql } from '@apollo/client';

export const GET_CELLS_BY_SPREADSHEET_ID = gql(/* GraphQL */ `
    query getCellsBySpreadsheetId {
        getCellsBySpreadsheetId(spreadsheetId: "1") {
            spreadsheet { name }
            id
            rowIndex
            columnIndex
            rawValue
            computedValue
        }
    }
`);

export const GET_SPREADSHEET = gql(/* GraphQL */ `
    query getSpreadsheet {
        getSpreadsheet(id: "1") {
            name
            rowCount
            columnCount
        }
    }
`);

export const UPDATE_CELL = gql`
    mutation UpdateCell($cellId: String!, $rawValue: String!) {
        updateCell(id: $cellId, input: { rawValue: $rawValue }) {
            id
            rawValue
        }
    }
`;

export const UPDATE_CELL_BY_SPREADSHEET_ID_COLUMN_AND_ROW = gql`
    mutation UpdateCellBySpreadsheetIdColumnAndRow($spreadsheetId: String!, $columnIndex: Int!, $rowIndex: Int!, $rawValue: String!) {
        updateCellBySpreadsheetIdColumnAndRow(spreadsheetId: $spreadsheetId, columnIndex: $columnIndex, rowIndex: $rowIndex, input: { rawValue: $rawValue }) {
            id
            rawValue
        }
    }
`;