/* jshint -W097,-W117 */
'use strict';

require("babel-polyfill");

var React = require('react');
var ReactDOM = require('react-dom');

var _ = require('underscore');

var utils = require('./utils.js');
var api = require('./api.js');
var action = require('./action.js');
var view = require('./view.js');

var ConnectionWindow = require('./ConnectionWindow.jsx');
var Sidebar = require('./Sidebar.jsx');
var AlertBar = require('./AlertBar.jsx');
var MainContainer = require('./MainContainer.jsx');
var SpinnerCircle = require('./Spinners.jsx').Circle;

const minSidebarDx = 128;
const maxSidebarDx = 128*3;

function isCreateOrDropQuery(query) {
  return query.match(/(create|drop) table/i);
}

function  isSelectQuery(query) {
  return query.match(/select/i);
}

class App extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleAlertBar = this.handleAlertBar.bind(this);
    this.handleCloseAlertBar = this.handleCloseAlertBar.bind(this);
    this.handleDidConnect = this.handleDidConnect.bind(this);
    this.handleDisconnectDatabase = this.handleDisconnectDatabase.bind(this);
    this.handleExecuteQuery = this.handleExecuteQuery.bind(this);
    this.handleExplainQuery = this.handleExplainQuery.bind(this);
    this.handleTableSelected = this.handleTableSelected.bind(this);
    this.handleViewSelected = this.handleViewSelected.bind(this);
    this.handleToggleSpinner = this.handleToggleSpinner.bind(this);
    this.handleResetPagination = this.handleResetPagination.bind(this);
    this.handleSelectedCellPosition = this.handleSelectedCellPosition.bind(this);
    this.handleEditedCells = this.handleEditedCells.bind(this);
    this.getQueryAsyncStatus = this.getQueryAsyncStatus.bind(this);
    this.getQueryAsyncData = this.getQueryAsyncData.bind(this);
    this.handleQueryAsync = this.handleQueryAsync.bind(this);
    this.handleQuerySync = this.handleQuerySync.bind(this);

    this.onMouseDown = this.onMouseDown.bind(this);
    this.onMouseMove = this.onMouseMove.bind(this);
    this.onMouseUp = this.onMouseUp.bind(this);

    this.state = {
      selectedView: view.SQLQuery,
      connectionId: gUserInfo ? gUserInfo.ConnectionID : 0,
      connected: gUserInfo ? gUserInfo.ConnectionID !== 0 : false,

      databaseName: "No Database Selected",

      queryIdInProgress: null,
      queryStatus: null,

      tables: null,
      tableStructures: {},
      selectedTable: "",
      selectedTableInfo: null,
      results: null,
      resetPagination: false,

      selectedCellPosition: {rowId: -1, colId: -1},
      editedCells: {},

      errorMessage: "",
      errorVisible: false,

      dragging: false,
      dragBarPosition: 250,

      spinnerVisible: 0,

      capabilities: {},
    };
  }

  getAllTablesStructures(connId, tables) {
    // We can do this by having the query get all data at once but its harder
    var self = this;
    _.each(tables, function(table) {
      api.getTableStructure(connId, table, function(tableStructureData) {
        var tempTableStructures = self.state.tableStructures;
        tempTableStructures[table] = tableStructureData;
        self.setState({
          tableStructures: tempTableStructures
        });

        console.log("All Table structrues, ", self.state.tableStructures);
      });
    });
  }

  handleDidConnect(connectionStr, connectionId, databaseName, capabilities) {
    this.setState({
      connected: true,
      connectionId: connectionId,
      databaseName: databaseName,
      capabilities: capabilities
    });
    var self = this;
    var connId = this.state.connectionId;
    api.getTables(connId, function(data) {
      self.setState({
        tables: data,
      });

      self.getAllTablesStructures(connId, data);
    });
  }

  onDragStart(e) {
    console.log("onDragStart");
  }

  componentDidUpdate(props, state) {
    if (this.state.dragging && !state.dragging) {
      document.addEventListener('mousemove', this.onMouseMove);
      document.addEventListener('mouseup', this.onMouseUp);
    } else if (!this.state.dragging && state.dragging) {
      document.removeEventListener('mousemove', this.onMouseMove);
      document.removeEventListener('mouseup', this.onMouseUp);
    }
  }

  onMouseDown(e) {
    // only left mouse button
    if (e.button !== 0) return;
    this.setState({
      dragging: true,
    });
    e.stopPropagation();
    e.preventDefault();
  }

  onMouseUp(e) {
    this.setState({
      dragging: false,
    });
    e.stopPropagation();
    e.preventDefault();
  }

  onMouseMove(e) {
    if (!this.state.dragging) return;
    if ((e.pageX < minSidebarDx) || (e.pageX > maxSidebarDx)) {
      return;
    }
    this.setState({
      dragBarPosition: e.pageX,
    });
    e.stopPropagation();
    e.preventDefault();
  }

  handleTableSelected(table) {
    console.log("handleTableSelected: table: ", table);

    this.setState({
      selectedTable: table,
    });

    // must delay otherwise this.state.selectedTable will not be visible yet
    // in handleViewSelected
    var self = this;
    setTimeout(function() {
      self.handleViewSelected(view.SQLQuery);
    }, 200);

    var connId = this.state.connectionId;
    api.getTableInfo(connId, table, function(data) {
      self.setState({
        selectedTableInfo: data,
      });
    });
  }

  getTableContent() {
    const table = this.state.selectedTable;
    const query = `SELECT * FROM ${table};`;
    this.handleExecuteQuery(query);
  }

  getTableStructure() {
    var self = this;
    var connId = this.state.connectionId;
    var selectedTable = this.state.selectedTable;
    api.getTableStructure(connId, selectedTable, function(data) {
      console.log("getTableStructure: ", data);
      self.setState({
        results: data,
        selectedCellPosition: {rowId: -1, colId: -1},
        editedCells: {},
      });
    });
  }

  getTableIndexes() {
    var self = this;
    var connId = this.state.connectionId;
    var selectedTable = this.state.selectedTable;
    api.getTableIndexes(connId, selectedTable, function(data) {
      console.log("getTableIndexes: ", data);
      self.setState({
        results: data,
        selectedCellPosition: {rowId: -1, colId: -1},
        editedCells: {},
      });
    });
  }

  getHistory() {
    var self = this;
    var connId = this.state.connectionId;
    api.getHistory(connId, function(data) {
      console.log("getHistory: ", data);
      self.setState({
        results: data,
        selectedCellPosition: {rowId: -1, colId: -1},
        editedCells: {},
      });
    });
  }

  /*
  getBookmarks() {
    var self = this;
    api.getBookmarks(function(data) {
      console.log("getBookmarks: ", data);
      self.setState({
        results: data
      });
    });
  }
  */

  handleViewSelected(viewName) {
    console.log("handleViewSelected: ", viewName);
    this.setState({
      selectedView: viewName
    });

    if (this.state.connectionId === 0) {
      console.log("handleViewSelected: not connected, connectionId: ", this.state.connectionId);
      return;
    }

    // those don't require table being selected
    switch (viewName) {
      case view.History:
        this.getHistory();
        break;
    }

    if (this.state.selectedTable === "") {
      //console.log("handleViewSelected: no selectedTable");
      return;
    }

    switch (viewName) {
      case view.SQLQuery:
        if (this.state.selectedTable == "") {
          this.setState({
            results: null,
          });
        } else {
          this.getTableContent();
        }
        return;
      case view.Structure:
        this.getTableStructure();
        break;
      case view.Indexes:
        this.getTableIndexes();
        break;
      default:
        console.log("handleViewSelected: unknown view: ", viewName);
    }
  }

  handleQuerySync(query) {
    console.log("handleQuerySync", query);
    var self = this;
    var connId = this.state.connectionId;
    api.executeQuery(connId, query, function(data) {
      self.setState({
        results: data,
        resetPagination: true,
        selectedCellPosition: {rowId: -1, colId: -1},
        editedCells: {},
      });

      // refresh tables list if table was added or removed
      if (isCreateOrDropQuery(query)) {
        api.getTables(self.state.connectionId, function(data) {
          self.setState({
            tables: data,
          });
        });
      }
    });
  }

  getQueryAsyncData() {
    const queryId = this.state.queryIdInProgress;
    console.log(`getQueryAsyncData: queryId={queryId}`);
    if (queryId == "") {
      console.log("no async query in progress");
      return;
    }
    const count = this.state.queryStatus.rows_count;
    if (count == 0) {
      this.setState({
        results: null,
        resetPagination: true,
        selectedCellPosition: {rowId: -1, colId: -1},
        editedCells: {}
      });
      return;
    }
    const connId = this.state.connectionId;
    const start = 0;
    const columns = this.state.queryStatus.columns;
    api.queryAsyncData(connId, queryId, start, count, (data) => {
      const results = {
        columns: columns,
        rows: data.rows
      };
      this.setState({
        results: results,
        spinnerVisible: false,
        resetPagination: true,
        selectedCellPosition: {rowId: -1, colId: -1},
        editedCells: {},
      });
    });
  }

  getQueryAsyncStatus() {
    const queryId = this.state.queryIdInProgress;
    console.log(`getQueryAsyncStatus: queryId={queryId}`);
    if (queryId == "") {
      console.log("no async query in progress");
      return;
    }
    const connId = this.state.connectionId;
    api.queryAsyncStatus(connId, queryId, (data) => {
      const queryStatus = data; 
      this.setState({
        queryStatus: queryStatus,
        spinnerVisible: !queryStatus.finished,
      });
      // repeat until async query finishes
      if (!queryStatus.finished) {
        setTimeout(this.getQueryAsyncStatus, 1000);
      } else {
        this.getQueryAsyncData();
      }
    });
  }

  handleQueryAsync(query) {
    console.log("handleQueryAsync", query);
    const connId = this.state.connectionId;
    api.queryAsync(connId, query, (data) => {
      this.setState({
        spinnerVisible: true,
        queryIdInProgress: data.query_id,
        // TODO: not sure if should reset the data right away
        // maybe only after received some data or an error message
        resetPagination: true,
        selectedCellPosition: {rowId: -1, colId: -1},
        editedCells: {},
      });
      setTimeout(this.getQueryAsyncStatus, 1000);
    });
  }

  handleExecuteQuery(query) {
    console.log("handleExecuteQuery", query);
    query = query.trim();
    if (isSelectQuery(query)) {
      this.handleQueryAsync(query);
    } else {
      this.handleQuerySync(query);
    }
  }

  handleExplainQuery(query) {
    console.log("handleExplainQuery", query);
    var self = this;
    var connId = this.state.connectionId;
    api.explainQuery(connId, query, function(data) {
      self.setState({
        selectedView: view.SQLQuery,
        results: data,
        resetPagination: true,
        selectedCellPosition: {rowId: -1, colId: -1},
        editedCells: {},
      });
    });
  }

  adHocTest() {
    var cid1 = action.onViewSelected(this.handleViewSelected);
    var cid2 = action.onViewSelected(this.handleViewSelected);
    action.offViewSelected(cid2);
    action.offViewSelected(cid1);
    action.offViewSelected(18);
  }

  handleDisconnectDatabase() {
    var self = this;
    api.disconnect(this.state.connectionId, function(data) {
      console.log("disconnect");

      self.setState({
        connectionId: 0,
        connected: false,
        tables: null,
        selectedTable: "",
        selectedTableInfo: null,
        results: null,
        errorMessage: "",
        errorVisible: false,
        resetPagination: false,
        selectedCellPosition: {rowId: -1, colId: -1},
        editedCells: {},
        spinnerVisible: 0,
      });
    });
  }

  handleAlertBar(message) {
    console.log("Create Alert Bar");

    this.setState({
      errorVisible: true,
      errorMessage: message,
    });

    // Dissmiss AlertBar with timer
    // var self = this;
    // setTimeout(function() {
    //   self.handleCloseAlertBar()
    // }, 5000);
  }

  handleCloseAlertBar() {
    this.setState({errorVisible: false});
  }

  handleToggleSpinner(toggle) {
    this.setState({spinnerVisible: toggle});
  }

  handleResetPagination(toggle) {
    this.setState({resetPagination: toggle});
  }

  handleSelectedCellPosition(newPosition) {
    this.setState({selectedCellPosition: newPosition});
  }

  handleEditedCells(newCells) {
    this.setState({editedCells: newCells});
  }

  componentWillMount() {
    //this.adHocTest();

    this.cidViewSelected = action.onViewSelected(this.handleViewSelected);
    this.cidTableSelected = action.onTableSelected(this.handleTableSelected);
    this.cidExecuteQuery = action.onExecuteQuery(this.handleExecuteQuery);
    this.cidExplainQuery = action.onExplainQuery(this.handleExplainQuery);
    this.cidDisconnectDatabase = action.onDisconnectDatabase(this.handleDisconnectDatabase);
    this.cidAlertBar = action.onAlertBar(this.handleAlertBar);
    this.cidSpinner = action.onSpinner(this.handleToggleSpinner);
    this.cidResetPagination = action.onResetPagination(this.handleResetPagination);
    this.cidSelectedCellPosition = action.onSelectedCellPosition(this.handleSelectedCellPosition);
    this.cidEditedCells = action.onEditedCells(this.handleEditedCells);

    var connId = this.state.connectionId;
    var self = this;
    if (connId !== 0) {
      api.getTables(connId, function(data) {
        self.setState({
          tables: data,
        });

       self.getAllTablesStructures(connId, data);
      });

      api.getConnectionInfo(connId, function(data) {
        self.setState({
          databaseName: _.filter(data.rows, function (el) { return (el[0] == "current_database"); })[0][1],
        });
      });
    }
  }

  componentWillUnmount() {
    action.offViewSelected(this.cidViewSelected);
    action.offTableSelected(this.cidTableSelected);
    action.offExecuteQuery(this.cidExecuteQuery);
    action.offExplainQuery(this.cidExplainQuery);
    action.offDisconnectDatabase(this.cidDisconnectDatabase);
    action.offAlertBar(this.cidAlertBar);
    action.offSpinner(this.cidSpinner);
    action.offResetPagination(this.cidResetPagination);
    action.offSelectedCellPosition(this.cidSelectedCellPosition);
    action.offEditedCells(this.cidEditedCells);
  }

  render() {
    var spinnerStyle = {
      position: 'fixed',
      top: '50%',
      left: '50%',
      zIndex: '5',
    };

    if (!this.state.connected) {
      return (
        <div>
          <div onClick={this.handleCloseAlertBar} >
            { this.state.errorVisible ? <AlertBar errorMessage={this.state.errorMessage}/> : null }
          </div>

          <SpinnerCircle style={spinnerStyle} visible={this.state.spinnerVisible} />
          <ConnectionWindow onDidConnect={this.handleDidConnect} />
        </div>
      );
    }

    var dragBarStyle = { left: this.state.dragBarPosition + 'px' };

    return (
      <div>
        <div onClick={this.handleCloseAlertBar} >
          { this.state.errorVisible ? <AlertBar errorMessage={this.state.errorMessage}/> : null }
        </div>

        <div>
          <Sidebar
            connectionId={this.state.connectionId}
            tables={this.state.tables}
            selectedTable={this.state.selectedTable}
            selectedTableInfo={this.state.selectedTableInfo}
            databaseName={this.state.databaseName}
            dragBarPosition={this.state.dragBarPosition} />

          <div id="side-dragbar"
            style={dragBarStyle}
            onMouseDown={this.onMouseDown}
            onMouseMove={this.onMouseMove}
            onMouseUp={this.onMouseUp}>
          </div>

          <MainContainer
            spinnerVisible={this.state.spinnerVisible}
            results={this.state.results}
            supportsExplain={this.state.capabilities.HasAnalyze}
            dragBarPosition={this.state.dragBarPosition}
            selectedView={this.state.selectedView}
            resetPagination={this.state.resetPagination}
            tableStructures={this.state.tableStructures}
            selectedTable={this.state.selectedTable}
            selectedCellPosition={this.state.selectedCellPosition}
            editedCells={this.state.editedCells} />
        </div>

      </div>
    );
  }
}

function appStart() {
  ReactDOM.render(
    <App/>,
    document.getElementById('main')
  );
}

window.appStart = appStart;
