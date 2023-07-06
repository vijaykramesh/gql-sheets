/* eslint-disable */
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>;
};
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>;
};
export type MakeEmpty<
  T extends { [key: string]: unknown },
  K extends keyof T,
> = { [_ in K]?: never };
export type Incremental<T> =
  | T
  | {
      [P in keyof T]?: P extends " $fragmentName" | "__typename" ? T[P] : never;
    };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string };
  String: { input: string; output: string };
  Boolean: { input: boolean; output: boolean };
  Int: { input: number; output: number };
  Float: { input: number; output: number };
};

export type Cell = {
  __typename?: "Cell";
  columnIndex: Scalars["Int"]["output"];
  computedValue?: Maybe<Scalars["String"]["output"]>;
  id: Scalars["String"]["output"];
  rawValue: Scalars["String"]["output"];
  rowIndex: Scalars["Int"]["output"];
  spreadsheet: Spreadsheet;
  version: Scalars["String"]["output"];
};

export type Mutation = {
  __typename?: "Mutation";
  createCell: Cell;
  createSpreadsheet: Spreadsheet;
  revertSpreadsheet: Spreadsheet;
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

export type MutationRevertSpreadsheetArgs = {
  id: Scalars["String"]["input"];
  version: Scalars["String"]["input"];
};

export type MutationUpdateCellArgs = {
  id: Scalars["String"]["input"];
  input: UpdateCell;
};

export type MutationUpdateCellBySpreadsheetIdColumnAndRowArgs = {
  columnIndex: Scalars["Int"]["input"];
  input: UpdateCell;
  rowIndex: Scalars["Int"]["input"];
  spreadsheetId: Scalars["String"]["input"];
};

export type MutationUpdateSpreadsheetArgs = {
  id: Scalars["String"]["input"];
  input: UpdateSpreadsheet;
};

export type NewCell = {
  columnIndex: Scalars["Int"]["input"];
  rawValue: Scalars["String"]["input"];
  rowIndex: Scalars["Int"]["input"];
  spreadsheetId: Scalars["String"]["input"];
};

export type NewSpreadsheet = {
  columnCount: Scalars["Int"]["input"];
  name: Scalars["String"]["input"];
  rowCount: Scalars["Int"]["input"];
};

export type Query = {
  __typename?: "Query";
  cells: Array<Cell>;
  getCell: Cell;
  getCellsBySpreadsheetId: Array<Cell>;
  getSpreadsheet: Spreadsheet;
  getVersions: Array<Version>;
  spreadsheets: Array<Spreadsheet>;
};

export type QueryGetCellArgs = {
  id: Scalars["String"]["input"];
};

export type QueryGetCellsBySpreadsheetIdArgs = {
  spreadsheetId: Scalars["String"]["input"];
};

export type QueryGetSpreadsheetArgs = {
  id: Scalars["String"]["input"];
};

export type QueryGetVersionsArgs = {
  id: Scalars["String"]["input"];
};

export type Spreadsheet = {
  __typename?: "Spreadsheet";
  columnCount: Scalars["Int"]["output"];
  id: Scalars["String"]["output"];
  name: Scalars["String"]["output"];
  rowCount: Scalars["Int"]["output"];
};

export type Subscription = {
  __typename?: "Subscription";
  getCellsBySpreadsheetId: Array<Cell>;
  getVersions: Array<Version>;
};

export type SubscriptionGetCellsBySpreadsheetIdArgs = {
  spreadsheetId: Scalars["String"]["input"];
};

export type SubscriptionGetVersionsArgs = {
  id: Scalars["String"]["input"];
};

export type UpdateCell = {
  rawValue: Scalars["String"]["input"];
};

export type UpdateSpreadsheet = {
  columnCount?: InputMaybe<Scalars["Int"]["input"]>;
  name?: InputMaybe<Scalars["String"]["input"]>;
  rowCount?: InputMaybe<Scalars["Int"]["input"]>;
};

export type Version = {
  __typename?: "Version";
  version: Scalars["String"]["output"];
};
