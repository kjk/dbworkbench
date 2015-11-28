/* jshint -W097,-W117 */
'use strict';

var _ = require('underscore');

var utils = require('./utils.js');
var api = require('./api.js');
var action = require('./action.js');
var view = require('./view.js');

var ConnectionWindow = require('./ConnectionWindow.jsx');
var DbNav = require('./DbNav.jsx');
var Sidebar = require('./Sidebar.jsx');
var Input = require('./Input.jsx');
var Output = require('./Output.jsx');

var App = React.createClass({
  getInitialState: function() {
    return {
      selectedView: view.SQLQuery,
      connectionId: gUserInfo ? gUserInfo.ConnectionID : 0,
      connected: gUserInfo ? gUserInfo.ConnectionID !== 0 : false,
      databaseName: "fixme: db name", // TODO: get database name
      tables: null,
      selectedTable: "",
      selectedTableInfo: null,
      results: null,
    };
  },

  handleDidConnect: function(connectionStr, connectionId, databaseName) {
    this.setState({
      connected: true,
      connectionId: connectionId,
      databaseName: databaseName
    });
    var self = this;
    var connId = this.state.connectionId;
    api.getTables(connId, function(data) {
      self.setState({
        tables: data,
      });
    });
  },

  renderInput: function() {
    if (this.state.selectedView === view.SQLQuery) {
      return <Input />;
    }
  },

  handleTableSelected: function(table) {
    console.log("handleTableSelected: table: ", table);
    this.setState({
      selectedTable: table,
    });

    // must delay otherwise this.state.selectedTable will not be visible yet
    // in handleViewSelected
    setTimeout(function() {
      action.viewSelected(view.Content);
    }, 200);

    var self = this;
    var connId = this.state.connectionId;
    api.getTableInfo(connId, table, function(data) {
      //console.log("handleTableSelected: tableInfo: ", data);
      self.setState({
        selectedTableInfo: data,
      });
    });
  },

  getTableContent: function() {
    var sortColumn = null;
    var sortOrder = null;
    var params = { limit: 100, sort_column: sortColumn, sort_order: sortOrder };

    var self = this;
    var connId = this.state.connectionId;
    var selectedTable = this.state.selectedTable;
    api.getTableRows(connId, selectedTable, params, function(data) {
      console.log("getTableContent: ", data);
      self.setState({
        results: data
      });
    });
  },

  getTableStructure: function() {
    var self = this;
    var connId = this.state.connectionId;
    var selectedTable = this.state.selectedTable;
    api.getTableStructure(connId, selectedTable, function(data) {
      console.log("getTableStructure: ", data);
      self.setState({
        results: data
      });
    });
  },

  getTableIndexes: function() {
    var self = this;
    var connId = this.state.connectionId;
    var selectedTable = this.state.selectedTable;
    api.getTableIndexes(connId, selectedTable, function(data) {
      console.log("getTableIndexes: ", data);
      self.setState({
        results: data
      });
    });
  },

  getHistory: function() {
    var self = this;
    var connId = this.state.connectionId;
    api.getHistory(connId, function(data) {
      console.log("getHistory: ", data);
      self.setState({
        results: data
      });
    });
  },

  getBookmarks: function() {
    var self = this;
    api.getBookmarks(function(data) {
      console.log("getBookmarks: ", data);
      self.setState({
        results: data
      });
    });
  },

  getActivity: function() {
    var self = this;
    var connId = this.state.connectionId;
    api.getActivity(connId, function(data) {
      console.log("getActivity: ", data);
      self.setState({
        results: data
      });
    });
  },

  getConnectionInfo: function() {
    var self = this;
    var connId = this.state.connectionId;
    api.getConnectionInfo(connId, function(data) {
      self.setState({
        results: data
      });
    });
  },

  handleViewSelected: function(viewName) {
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
      case view.SQLQuery:
        this.setState({
          results: null,
        });
        return;
      case view.History:
        this.getHistory();
        break;
      case view.Connection:
        this.getConnectionInfo();
        return;
      case view.Activity:
        this.getActivity();
        return;
    }

    if (this.state.selectedTable === "") {
      //console.log("handleViewSelected: no selectedTable");
      return;
    }

    switch (viewName) {
      case view.Content:
        this.getTableContent();
        break;
      case view.Structure:
        this.getTableStructure();
        break;
      case view.Indexes:
        this.getTableIndexes();
        break;
      default:
        console.log("handleViewSelected: unknown view: ", viewName);
    }
  },

  handleExecuteQuery: function(query) {
    console.log("handleExecuteQuery", query);
    var self = this;
    var connId = this.state.connectionId;
    api.executeQuery(connId, query, function(data) {
      self.setState({
        results: data
      });

      // refresh tables list if table was added or removed
      var re = /(create|drop) table/i;
      if (query.match(re)) {
        api.getTables(self.state.connectionId, function(data) {
          self.setState({
            tables: data,
          });
        });
      }
    });
  },

  handleExplainQuery: function(query) {
    console.log("handleExplainQuery", query);
    var self = this;
    var connId = this.state.connectionId;
    api.explainQuery(connId, query, function(data) {
      self.setState({
        selectedView: view.SQLQuery,
        results: data
      });
    });
  },

  adHocTest: function() {
    var cid1 = action.onViewSelected(this.handleViewSelected);
    var cid2 = action.onViewSelected(this.handleViewSelected);
    action.offViewSelected(cid2);
    action.offViewSelected(cid1);
    action.offViewSelected(18);
  },

  componentDidMount: function() {
    //this.adHocTest();

    this.cidViewSelected = action.onViewSelected(this.handleViewSelected);
    this.cidTableSelected = action.onTableSelected(this.handleTableSelected);
    this.cidExecuteQuery = action.onExecuteQuery(this.handleExecuteQuery);
    this.cidExplainQuery = action.onExplainQuery(this.handleExplainQuery);

    var connId = this.state.connectionId;
    var self = this;
    if (connId !== 0) {
      api.getTables(connId, function(data) {
        self.setState({
          tables: data,
        });
      });

      api.getConnectionInfo(connId, function(data) {
        self.setState({
          databaseName: _.filter(data.rows, function (el) { return (el[0] == "current_database"); })[0][1],
        });
      });
    }
  },

  componentDidUnmount: function() {
    action.offViewSelected(this.cidViewSelected);
    action.offTableSelected(this.cidTableSelected);
    action.offExecuteQuery(this.cidExecuteQuery);
    action.offExplainQuery(this.cidExplainQuery);
  },

  render: function() {
    if (!this.state.connected) {
      return <ConnectionWindow onDidConnect={this.handleDidConnect} />;
    }

    // when showing sql query, results are below editor window
    var notFull = (this.state.selectedView === view.SQLQuery);
    return (
      <div>
        <DbNav view={this.state.selectedView}/>
        <Sidebar
          tables={this.state.tables}
          selectedTable={this.state.selectedTable}
          selectedTableInfo={this.state.selectedTableInfo}
          databaseName={this.state.databaseName}
        />
        <div id="body">
          {this.renderInput()}
          <Output selectedView={this.state.selectedView} results={this.state.results} notFull={notFull}/>
        </div>
      </div>
    );
  }
});

function appStart() {
  React.render(
    <App/>,
    document.getElementById('main')
  );
}

window.appStart = appStart;
