/* jshint -W097,-W117 */
'use strict';

var utils = require('./utils.js');
var api = require('./api.js');
var action = require('./action.js');
var view = require('./view.js');

var ConnectionWindow = require('./ConnectionWindow.jsx');
var TopNav = require('./TopNav.jsx').TopNav;
var ViewSQLQuery = require('./TopNav.jsx').ViewSQLQuery;
var Sidebar = require('./Sidebar.jsx');
var Input = require('./Input.jsx');
var Output = require('./Output.jsx');

var App = React.createClass({
  getInitialState: function() {
    return {
      selectedView: view.SQLQuery,
      connectionId: -1,
      connected: false,
      databaseName: "",
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
    api.getTables(function(data) {
      self.setState({
        tables: data,
      });
    });
  },

  renderInput: function() {
    if (this.selectedView === view.SQLQuery) {
      return <Input />;
    }
  },

  handleTableSelected: function(table) {
    console.log("handleTableSelected: table: ", table);
    this.setState({
      selectedTable: table,
      //selectedTableInfo: null,
      //results: null
    });

    // must delay otherwise this.state.selectedTable will not be visible yet
    // in handleViewSelected
    setTimeout(function() {
      action.viewSelected(view.Content);
    }, 200);

    var self = this;
    api.call("get", "/tables/" + table + "/info", {}, function(data) {
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
    api.getTableRows(this.state.selectedTable, params, function(data) {
      console.log("getTableContent: ", data);
      self.setState({
        results: data
      });
    });
  },

  getTableStructure: function() {
    var self = this;
    api.getTableStructure(this.state.selectedTable, function(data) {
      console.log("getTableStructure: ", data);
      self.setState({
        results: data
      });
    });
  },

  getTableIndexes: function() {
    var self = this;
    api.getTableIndexes(this.state.selectedTable, function(data) {
      console.log("getTableIndexes: ", data);
      self.setState({
        results: data
      });
    });
  },

  getHistory: function() {
    var self = this;
    api.getHistory(function(data) {
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
    api.getActivity(function(data) {
      console.log("getActivity: ", data);
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

    if (this.state.connectionId === -1) {
      console.log("handleViewSelected: not connected, connectionId: ", this.state.connectionId);
      return;
    }
    if (this.state.selectedTable === "") {
      console.log("handleViewSelected: no selectedTable");
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
      case view.SQLQuery:
        this.setState({
          results: null,
        });
        break;
      case view.History:
        this.getHistory();
        break;
      case view.Activity:
        this.getActivity();
        break;
      case view.Connection:
        // TODO: write me
        break;
      default:
        console.log("handleViewSelected: unknown view: ", viewName);
    }
  },

  componentDidMount: function() {
    action.onViewSelected(this.handleViewSelected);
    action.onTableSelected(this.handleTableSelected);
  },

  componentDidUnmount: function() {
    action.cancelOnViewSelected(this.handleViewSelected);
    action.cancelOnTableSelected(this.handleTableSelected);
  },

  render: function() {
    if (!this.state.connected) {
      return <ConnectionWindow onDidConnect={this.handleDidConnect} />;
    }

    return (
      <div>
        <TopNav view={this.state.selectedView}/>
        <Sidebar
          tables={this.state.tables}
          selectedTable={this.state.selectedTable}
          selectedTableInfo={this.state.selectedTableInfo}
          databaseName={this.state.databaseName}
        />
        <div id="body">
          {this.renderInput()}
          <Output results={this.state.results}/>
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
