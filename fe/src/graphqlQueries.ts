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
            id
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

export const UPDATE_SPREADSHEET = gql`
    mutation UpdateSpreadsheet($id: String!, $rowCount: Int!, $columnCount: Int!) {
        updateSpreadsheet(id: $id, input: {  rowCount: $rowCount, columnCount: $columnCount }) {
            id
            name
            rowCount
            columnCount
        }
    }
`;

export const GET_CELLS_BY_SPREADSHEET_ID_SUBSCRIPTION = gql`
    subscription getCellsBySpreadsheetId($spreadsheetId: String!) {
        getCellsBySpreadsheetId(spreadsheetId: $spreadsheetId) {
            spreadsheet { name }
            id
            rowIndex
            columnIndex
            rawValue
            computedValue
        }
    }
`;

export const GET_VERSIONS_BY_SPREADSHEET_ID = gql(/* GraphQL */ `
    query getVersionsBySpreadsheetId($spreadsheetId: String!)      {
        getVersions(id: $spreadsheetId) {
            version
        }
    }
`);

export const REVERT_SPREADSHEET_TO_VERSION = gql`
    mutation revertSpreadsheet($spreadsheetId: String!, $version: String!) {
        revertSpreadsheet(id: $spreadsheetId, version: $version) {
            id
            name
            rowCount
            columnCount
        }
    }
`;