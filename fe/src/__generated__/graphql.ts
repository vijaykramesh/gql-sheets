/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export type Cell = {
  __typename?: 'Cell';
  columnIndex: Scalars['Int']['output'];
  computedValue?: Maybe<Scalars['String']['output']>;
  id: Scalars['String']['output'];
  rawValue: Scalars['String']['output'];
  rowIndex: Scalars['Int']['output'];
  spreadsheet: Spreadsheet;
};

export type Mutation = {
  __typename?: 'Mutation';
  createCell: Cell;
  createSpreadsheet: Spreadsheet;
  updateCell: Cell;
  updateCellBySpreadsheetIdColumnAndRow: Cell;
  updateSpreadsheet: Spreadsheet;
};


export type MutationCreateCellArgs = {
  input: NewCell;
};


export type MutationCreateSpreadsheetArgs = {
  input: NewSpreadsheet;
};


export type MutationUpdateCellArgs = {
  id: Scalars['String']['input'];
  input: UpdateCell;
};


export type MutationUpdateCellBySpreadsheetIdColumnAndRowArgs = {
  columnIndex: Scalars['Int']['input'];
  input: UpdateCell;
  rowIndex: Scalars['Int']['input'];
  spreadsheetId: Scalars['String']['input'];
};


export type MutationUpdateSpreadsheetArgs = {
  id: Scalars['String']['input'];
  input: NewSpreadsheet;
};

export type NewCell = {
  columnIndex: Scalars['Int']['input'];
  rawValue: Scalars['String']['input'];
  rowIndex: Scalars['Int']['input'];
  spreadsheetId: Scalars['String']['input'];
};

export type NewSpreadsheet = {
  columnCount: Scalars['Int']['input'];
  name: Scalars['String']['input'];
  rowCount: Scalars['Int']['input'];
};

export type Query = {
  __typename?: 'Query';
  cells: Array<Cell>;
  getCell: Cell;
  getCellsBySpreadsheetId: Array<Cell>;
  getSpreadsheet: Spreadsheet;
  spreadsheets: Array<Spreadsheet>;
};


export type QueryGetCellArgs = {
  id: Scalars['String']['input'];
};


export type QueryGetCellsBySpreadsheetIdArgs = {
  spreadsheetId: Scalars['String']['input'];
};


export type QueryGetSpreadsheetArgs = {
  id: Scalars['String']['input'];
};

export type Spreadsheet = {
  __typename?: 'Spreadsheet';
  columnCount: Scalars['Int']['output'];
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
  rowCount: Scalars['Int']['output'];
};

export type UpdateCell = {
  rawValue: Scalars['String']['input'];
};

export type GetCellsBySpreadsheetIdQueryVariables = Exact<{ [key: string]: never; }>;


export type GetCellsBySpreadsheetIdQuery = { __typename?: 'Query', getCellsBySpreadsheetId: Array<{ __typename?: 'Cell', id: string, rowIndex: number, columnIndex: number, rawValue: string, computedValue?: string | null, spreadsheet: { __typename?: 'Spreadsheet', name: string } }> };

export type GetSpreadsheetQueryVariables = Exact<{ [key: string]: never; }>;


export type GetSpreadsheetQuery = { __typename?: 'Query', getSpreadsheet: { __typename?: 'Spreadsheet', name: string, rowCount: number, columnCount: number } };

export type UpdateCellMutationVariables = Exact<{
  cellId: Scalars['String']['input'];
  rawValue: Scalars['String']['input'];
}>;


export type UpdateCellMutation = { __typename?: 'Mutation', updateCell: { __typename?: 'Cell', id: string, rawValue: string } };

export type UpdateCellBySpreadsheetIdColumnAndRowMutationVariables = Exact<{
  spreadsheetId: Scalars['String']['input'];
  columnIndex: Scalars['Int']['input'];
  rowIndex: Scalars['Int']['input'];
  rawValue: Scalars['String']['input'];
}>;


export type UpdateCellBySpreadsheetIdColumnAndRowMutation = { __typename?: 'Mutation', updateCellBySpreadsheetIdColumnAndRow: { __typename?: 'Cell', id: string, rawValue: string } };


export const GetCellsBySpreadsheetIdDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"getCellsBySpreadsheetId"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"getCellsBySpreadsheetId"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"spreadsheetId"},"value":{"kind":"StringValue","value":"1","block":false}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"spreadsheet"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}}]}},{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"rowIndex"}},{"kind":"Field","name":{"kind":"Name","value":"columnIndex"}},{"kind":"Field","name":{"kind":"Name","value":"rawValue"}},{"kind":"Field","name":{"kind":"Name","value":"computedValue"}}]}}]}}]} as unknown as DocumentNode<GetCellsBySpreadsheetIdQuery, GetCellsBySpreadsheetIdQueryVariables>;
export const GetSpreadsheetDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"getSpreadsheet"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"getSpreadsheet"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"StringValue","value":"1","block":false}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"rowCount"}},{"kind":"Field","name":{"kind":"Name","value":"columnCount"}}]}}]}}]} as unknown as DocumentNode<GetSpreadsheetQuery, GetSpreadsheetQueryVariables>;
export const UpdateCellDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateCell"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"cellId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"rawValue"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateCell"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"cellId"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"rawValue"},"value":{"kind":"Variable","name":{"kind":"Name","value":"rawValue"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"rawValue"}}]}}]}}]} as unknown as DocumentNode<UpdateCellMutation, UpdateCellMutationVariables>;
export const UpdateCellBySpreadsheetIdColumnAndRowDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateCellBySpreadsheetIdColumnAndRow"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"spreadsheetId"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"columnIndex"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"rowIndex"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"Int"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"rawValue"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateCellBySpreadsheetIdColumnAndRow"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"spreadsheetId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"spreadsheetId"}}},{"kind":"Argument","name":{"kind":"Name","value":"columnIndex"},"value":{"kind":"Variable","name":{"kind":"Name","value":"columnIndex"}}},{"kind":"Argument","name":{"kind":"Name","value":"rowIndex"},"value":{"kind":"Variable","name":{"kind":"Name","value":"rowIndex"}}},{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"ObjectValue","fields":[{"kind":"ObjectField","name":{"kind":"Name","value":"rawValue"},"value":{"kind":"Variable","name":{"kind":"Name","value":"rawValue"}}}]}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"rawValue"}}]}}]}}]} as unknown as DocumentNode<UpdateCellBySpreadsheetIdColumnAndRowMutation, UpdateCellBySpreadsheetIdColumnAndRowMutationVariables>;