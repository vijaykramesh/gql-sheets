import React, {
  useState,
  useEffect,
  useMemo,
  useCallback,
  FunctionComponent,
} from "react";
import { useQuery, useMutation, useSubscription } from "@apollo/client";
import { AgGridReact } from "ag-grid-react";
import {
  ColumnApi,
  GridReadyEvent,
  ModuleRegistry,
  CellEditingStartedEvent,
  CellEditingStoppedEvent,
  CellValueChangedEvent,
  ColDef,
} from "ag-grid-community";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  createTheme,
  ThemeProvider,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Typography,
  CircularProgress,
  Box,
} from "@mui/material";
import "ag-grid-community/styles/ag-grid.css";
import "ag-grid-community/styles/ag-theme-alpine.css";

import {
  GET_CELLS_BY_SPREADSHEET_ID,
  GET_CELLS_BY_SPREADSHEET_ID_SUBSCRIPTION,
  GET_SPREADSHEET,
  GET_VERSIONS_BY_SPREADSHEET_ID,
  GET_VERSIONS_SUBSCRIPTION,
  REVERT_SPREADSHEET_TO_VERSION,
  UPDATE_CELL_BY_SPREADSHEET_ID_COLUMN_AND_ROW,
  UPDATE_SPREADSHEET,
} from "./graphqlQueries";
import { Cell } from "./__generated__/graphql";
import { ClientSideRowModelModule } from "@ag-grid-community/client-side-row-model";
import "./App.css";
import { generateColumnFields, convertCellsToRowData } from "./utils";

interface RowData {
  rowIndex: string;
  rawValue: string;
  computedValue: string;
  editMode: number;
}

ModuleRegistry.registerModules([ClientSideRowModelModule]);

const App: FunctionComponent = (): React.ReactElement => {
  const { loading, data } = useQuery(GET_CELLS_BY_SPREADSHEET_ID);
  const { loading: loadingSpreadsheet, data: dataSpreadsheet } =
    useQuery(GET_SPREADSHEET);
  const [updateCell, { data: dataUpdate, loading: loadingUpdate, error }] =
    useMutation(UPDATE_CELL_BY_SPREADSHEET_ID_COLUMN_AND_ROW, {
      refetchQueries: [
        { query: GET_CELLS_BY_SPREADSHEET_ID },
        {
          query: GET_VERSIONS_BY_SPREADSHEET_ID,
          variables: { spreadsheetId: "1" },
        },
      ],
    });

  const [
    updateSpreadsheet,
    {
      data: dataUpdateSpreadsheet,
      loading: loadingUpdateSpreadsheet,
      error: errorUpdateSpreadsheet,
    },
  ] = useMutation(UPDATE_SPREADSHEET, {
    refetchQueries: [
      { query: GET_SPREADSHEET },
      { query: GET_CELLS_BY_SPREADSHEET_ID },
      {
        query: GET_VERSIONS_BY_SPREADSHEET_ID,
        variables: { spreadsheetId: "1" },
      },
    ],
  });

  const [
    revertVersion,
    {
      data: dataRevertVersion,
      loading: loadingRevertVersion,
      error: errorRevertVersion,
    },
  ] = useMutation(REVERT_SPREADSHEET_TO_VERSION, {
    refetchQueries: [
      { query: GET_SPREADSHEET },
      { query: GET_CELLS_BY_SPREADSHEET_ID },
      {
        query: GET_VERSIONS_BY_SPREADSHEET_ID,
        variables: { spreadsheetId: "1" },
      },
    ],
  });

  const { data: dataSubscription, loading: loadingSubscription } =
    useSubscription(
      GET_CELLS_BY_SPREADSHEET_ID_SUBSCRIPTION,

      { variables: { spreadsheetId: "1" } }, // todo get spreadsheet id from somewhere
    );

  const {
    data: dataVersionsSubscription,
    loading: loadingVersionsSubscription,
  } = useSubscription(
    GET_VERSIONS_SUBSCRIPTION,
    { variables: { spreadsheetId: "1" } }, // todo get spreadsheet id from somewhere
  );

  const containerStyle = useMemo(() => {
    return { width: "100%", height: "100%" };
  }, []);
  const gridStyle = useMemo(() => {
    return { height: "800px", width: "90%", padding: "50px" };
  }, []);
  const [editModeColumnAndRow, setEditModeColumnAndRow] = useState<{
    column: string;
    row: number;
  } | null>(null);
  const [rowData, setRowData] = useState<any[]>();
  const [columnDefs, setColumnDefs] = useState<ColDef[]>([
    {
      headerName: "",
      field: "rowIndex",
      width: 50,
      sortable: false,
      resizable: false,
    },
    { field: "A" },
    { field: "B" },
    { field: "C" },
  ]);

  // VERSION SELECTOR
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false);
  const [selectedVersion, setSelectedVersion] = useState("");

  const handleVersionChange = (event) => {
    const newVersion = event.target.value;
    setSelectedVersion(newVersion);
    setIsConfirmationOpen(true); // Open the confirmation modal
  };

  const handleConfirmation = () => {
    // Perform the desired action (e.g., console.log)
    revertVersion({
      variables: { spreadsheetId: "1", version: selectedVersion },
    });
    setIsConfirmationOpen(false); // Close the confirmation modal
    setSelectedVersion(""); // Reset the selected version
  };

  const handleCloseConfirmation = () => {
    setIsConfirmationOpen(false); // Close the confirmation modal
    setSelectedVersion(""); // Reset the selected version
  };

  const defaultColDef = useMemo<ColDef>(() => {
    return {
      editable: true,
      cellDataType: false,
    };
  }, [data]);

  const onGridReady = useCallback(
    (params: GridReadyEvent) => {
      if (data && dataSpreadsheet) {
        const preparedRowData = convertCellsToRowData(
          data.getCellsBySpreadsheetId,
          dataSpreadsheet.getSpreadsheet.rowCount,
          editModeColumnAndRow,
        );
        setRowData(preparedRowData);
      } else {
        console.log("data is null");
      }
      if (dataSpreadsheet) {
        const columnDefs = generateColumnFields(
          dataSpreadsheet.getSpreadsheet.columnCount,
        );
        const firstColumn = [
          {
            headerName: "",
            field: "rowIndex",
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
        console.log("dataSpreadsheet is null");
      }
    },
    [data, dataSpreadsheet, editModeColumnAndRow],
  );

  useEffect(() => {
    if (data && dataSpreadsheet) {
      if (editModeColumnAndRow == null) {
        const preparedRowData = convertCellsToRowData(
          data.getCellsBySpreadsheetId,
          dataSpreadsheet.getSpreadsheet.rowCount,
          editModeColumnAndRow,
        );
        setRowData(preparedRowData);
      }
    } else {
      console.log("data is null");
    }
  }, [data, dataSpreadsheet, editModeColumnAndRow]);

  useEffect(() => {
    if (dataSubscription && dataSpreadsheet && dataSpreadsheet.getSpreadsheet) {
      if (editModeColumnAndRow == null) {
        const preparedRowData = convertCellsToRowData(
          dataSubscription.getCellsBySpreadsheetId,
          dataSpreadsheet.getSpreadsheet.rowCount,
          editModeColumnAndRow,
        );
        setRowData(preparedRowData);
      }
    } else {
      console.log("data is null");
    }
  }, [dataSubscription, dataSpreadsheet, editModeColumnAndRow]);

  useEffect(() => {
    document.title = "gql-sheets";
  }, []);

  const onCellEditingStarted = (event: CellEditingStartedEvent) => {
    const { rowIndex, column } = event;
    const columnCode = column.getColId();
    setEditModeColumnAndRow({ column: "-1", row: -1 });
    if (event.data && event.data[columnCode]) {
      event.data[columnCode].editMode = 2;
      if (rowIndex !== null && rowIndex >= 0 && columnCode) {
        setEditModeColumnAndRow({ column: columnCode, row: rowIndex });
      }
    }
  };

  const onCellEditingStopped = (event: CellEditingStoppedEvent) => {
    const { rowIndex, column } = event;
    const columnCode = column.getColId();

    if (columnCode && event.data && event.data[columnCode]) {
      event.data[columnCode].editMode = 0;
    }
    setEditModeColumnAndRow(null);

    if (!event.valueChanged) {
      return;
    }
    // if they have edited a cell in a new row or column we have to update the sheet first
    // the UX here is cludgy, but it's a demo
    if (
      dataSpreadsheet.getSpreadsheet.rowCount < (rowIndex || 0) + 1 ||
      dataSpreadsheet.getSpreadsheet.columnCount < column.getInstanceId()
    ) {
      const newRowCount = Math.max(
        (rowIndex || 0) + 1,
        dataSpreadsheet.getSpreadsheet.rowCount,
      );
      const newColumnCount = Math.max(
        column.getInstanceId(),
        dataSpreadsheet.getSpreadsheet.columnCount,
      );
      updateSpreadsheet({
        variables: {
          id: dataSpreadsheet.getSpreadsheet.id,
          rowCount: newRowCount,
          columnCount: newColumnCount,
        },
      });

      if (newColumnCount > dataSpreadsheet.getSpreadsheet.columnCount) {
        const columnDefs = generateColumnFields(newColumnCount);
        const firstColumn = [
          {
            headerName: "",
            field: "rowIndex",
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
        rawValue: event.newValue,
      },
    });
  };

  const gridOptions = {
    onCellEditingStarted: onCellEditingStarted,
    onCellEditingStopped: onCellEditingStopped,
    columnHoverHighlight: true,
  };

  if (loading || loadingSpreadsheet || loadingVersionsSubscription) {
    return (
      <Box
        display="flex"
        alignItems="center"
        justifyContent="center"
        height="100vh" // Adjust the height as needed
      >
        <CircularProgress />
      </Box>
    );
  }

  if (!data && !dataSpreadsheet) {
    return (
      <Box
        display="flex"
        alignItems="center"
        justifyContent="center"
        height="100vh" // Adjust the height as needed
      >
        <Typography variant="body1">Not found</Typography>
      </Box>
    );
  }
  const theme = createTheme();
  return (
    <ThemeProvider theme={theme}>
      <div className="App" style={{ height: "100%" }}>
        <header
          className="App-header"
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            height: "70px",
          }}
        >
          <Typography
            variant="h6"
            style={{ color: "#FFF", paddingTop: "20px" }}
          >
            gql-sheets
          </Typography>
        </header>

        <div className="container">
          <div
            className="revert-section"
            style={{
              marginLeft: "auto",
              width: "300px",
              marginRight: "50px",
              padding: "20px",
            }}
          >
            <FormControl
              variant="outlined"
              style={{ width: 300 }}
              className="revert-select"
            >
              {/* eslint-disable-next-line no-extra-parens */}
              {dataVersionsSubscription.getVersions && (
                <Select
                  value={
                    selectedVersion ||
                    (dataVersionsSubscription.getVersions[
                      dataVersionsSubscription.getVersions.length - 1
                    ]
                      ? dataVersionsSubscription.getVersions[
                          dataVersionsSubscription.getVersions.length - 1
                        ].version
                      : Date.now() / 1000)
                  }
                  onChange={handleVersionChange}
                >
                  {dataVersionsSubscription.getVersions.map(
                    (version, index) => {
                      return (
                        <MenuItem key={index} value={version.version}>
                          {new Date(Number(version.version)).toLocaleString()}
                        </MenuItem>
                      );
                    },
                  )}
                </Select>
              )}
            </FormControl>
          </div>
          <div className="modal">
            <Dialog open={isConfirmationOpen} onClose={handleCloseConfirmation}>
              <DialogTitle>Confirmation</DialogTitle>
              <DialogContent>
                Are you sure you want to proceed? This will revert the
                spreadsheet to a previous version and lose all changes made
                since then.
              </DialogContent>
              <DialogActions>
                <Button onClick={handleCloseConfirmation} color="primary">
                  Cancel
                </Button>
                <Button onClick={handleConfirmation} color="primary" autoFocus>
                  OK
                </Button>
              </DialogActions>
            </Dialog>
          </div>

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
    </ThemeProvider>
  );
};

export default App;
