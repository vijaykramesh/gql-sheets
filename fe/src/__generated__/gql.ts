/* eslint-disable */
import * as types from './graphql';
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 */
const documents = {
    "\n\n  query getCellsBySpreadsheetId {\n    getCellsBySpreadsheetId(spreadsheetId: \"1\"){\n        spreadsheet { name }\n        id\n        rowIndex\n        columnIndex\n        rawValue\n        computedValue\n      }\n    }\n\n": types.GetCellsBySpreadsheetIdDocument,
    "\n    query getSpreadsheet {\n        getSpreadsheet(id: \"1\") {\n            name\n            rowCount\n            columnCount\n        }\n    }\n": types.GetSpreadsheetDocument,
    "\n    mutation UpdateCell($cellId: String!, $rawValue: String!) {\n        updateCell(id: $cellId, input: { rawValue: $rawValue }) {\n            id\n            rawValue\n        }\n    }\n": types.UpdateCellDocument,
    "\n    mutation UpdateCellBySpreadsheetIdColumnAndRow($spreadsheetId: String!, $columnIndex: Int!, $rowIndex: Int!, $rawValue: String!) {\n        updateCellBySpreadsheetIdColumnAndRow(spreadsheetId: $spreadsheetId, columnIndex: $columnIndex, rowIndex: $rowIndex, input: { rawValue: $rawValue }) {\n            id\n            rawValue\n        }\n    }\n": types.UpdateCellBySpreadsheetIdColumnAndRowDocument,
};

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = gql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function gql(source: string): unknown;

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n\n  query getCellsBySpreadsheetId {\n    getCellsBySpreadsheetId(spreadsheetId: \"1\"){\n        spreadsheet { name }\n        id\n        rowIndex\n        columnIndex\n        rawValue\n        computedValue\n      }\n    }\n\n"): (typeof documents)["\n\n  query getCellsBySpreadsheetId {\n    getCellsBySpreadsheetId(spreadsheetId: \"1\"){\n        spreadsheet { name }\n        id\n        rowIndex\n        columnIndex\n        rawValue\n        computedValue\n      }\n    }\n\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n    query getSpreadsheet {\n        getSpreadsheet(id: \"1\") {\n            name\n            rowCount\n            columnCount\n        }\n    }\n"): (typeof documents)["\n    query getSpreadsheet {\n        getSpreadsheet(id: \"1\") {\n            name\n            rowCount\n            columnCount\n        }\n    }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n    mutation UpdateCell($cellId: String!, $rawValue: String!) {\n        updateCell(id: $cellId, input: { rawValue: $rawValue }) {\n            id\n            rawValue\n        }\n    }\n"): (typeof documents)["\n    mutation UpdateCell($cellId: String!, $rawValue: String!) {\n        updateCell(id: $cellId, input: { rawValue: $rawValue }) {\n            id\n            rawValue\n        }\n    }\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\n    mutation UpdateCellBySpreadsheetIdColumnAndRow($spreadsheetId: String!, $columnIndex: Int!, $rowIndex: Int!, $rawValue: String!) {\n        updateCellBySpreadsheetIdColumnAndRow(spreadsheetId: $spreadsheetId, columnIndex: $columnIndex, rowIndex: $rowIndex, input: { rawValue: $rawValue }) {\n            id\n            rawValue\n        }\n    }\n"): (typeof documents)["\n    mutation UpdateCellBySpreadsheetIdColumnAndRow($spreadsheetId: String!, $columnIndex: Int!, $rowIndex: Int!, $rawValue: String!) {\n        updateCellBySpreadsheetIdColumnAndRow(spreadsheetId: $spreadsheetId, columnIndex: $columnIndex, rowIndex: $rowIndex, input: { rawValue: $rawValue }) {\n            id\n            rawValue\n        }\n    }\n"];

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;