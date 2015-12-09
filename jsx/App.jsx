/* jshint -W097,-W117 */
'use strict';

var React = require('react');
var ReactDOM = require('react-dom');
var Spinner = require('react-spinkit');

var _ = require('underscore');

var utils = require('./utils.js');
var api = require('./api.js');
var action = require('./action.js');
var view = require('./view.js');

var ConnectionWindow = require('./ConnectionWindow.jsx');
var Sidebar = require('./Sidebar.jsx');
var AlertBar = require('./AlertBar.jsx');
var MainContainer = require('./MainContainer.jsx');

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
    this.onMouseDown = this.onMouseDown.bind(this);
    this.onMouseMove = this.onMouseMove.bind(this);
    this.onMouseUp = this.onMouseUp.bind(this);

    this.state = {
      selectedView: view.SQLQuery,
      connectionId: gUserInfo ? gUserInfo.ConnectionID : 0,
      connected: gUserInfo ? gUserInfo.ConnectionID !== 0 : false,

      databaseName: "fixme: db name", // TODO: get database name
      tables: null,
      selectedTable: "",
      selectedTableInfo: null,
      results: null,

      errorMessage: "",
      errorVisible: false,

      dragging: false,
      dragBarPosition: 250,

      spinnerVisible: false,
    };
  }

  handleDidConnect(connectionStr, connectionId, databaseName) {
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
  }

  onDragStart(e) {
    console.log("onDragStart");

  }

  componentDidUpdate(props, state) {
    if (this.state.dragging && !state.dragging) {
      document.addEventListener('mousemove', this.onMouseMove)
      document.addEventListener('mouseup', this.onMouseUp)
    } else if (!this.state.dragging && state.dragging) {
      document.removeEventListener('mousemove', this.onMouseMove)
      document.removeEventListener('mouseup', this.onMouseUp)
    }
  }

  onMouseDown(e) {
    // only left mouse button
    if (e.button !== 0) return;
    this.setState({
      dragging: true,
    })
    e.stopPropagation()
    e.preventDefault()
  }

  onMouseUp(e) {
    this.setState({
      dragging: false,
    })
    e.stopPropagation()
    e.preventDefault()
  }

  onMouseMove(e) {
    if (!this.state.dragging) return;
    this.setState({
      dragBarPosition: e.pageX,
    });
    e.stopPropagation()
    e.preventDefault()
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
      self.handleViewSelected(view.Content);
    }, 200);

    var connId = this.state.connectionId;
    api.getTableInfo(connId, table, function(data) {
      self.setState({
        selectedTableInfo: data,
      });
    });
  }

  getTableContent() {
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
  }

  getTableStructure() {
    var self = this;
    var connId = this.state.connectionId;
    var selectedTable = this.state.selectedTable;
    api.getTableStructure(connId, selectedTable, function(data) {
      console.log("getTableStructure: ", data);
      self.setState({
        results: data
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
        results: data
      });
    });
  }

  getHistory() {
    var self = this;
    var connId = this.state.connectionId;
    api.getHistory(connId, function(data) {
      console.log("getHistory: ", data);
      self.setState({
        results: data
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
      case view.SQLQuery:
        this.setState({
          results: null,
        });
        return;
      case view.History:
        this.getHistory();
        break;
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
  }

  handleExecuteQuery(query) {
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
  }

  handleExplainQuery(query) {
    console.log("handleExplainQuery", query);
    var self = this;
    var connId = this.state.connectionId;
    api.explainQuery(connId, query, function(data) {
      self.setState({
        selectedView: view.SQLQuery,
        results: data
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
            connected: false
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

  componentWillMount() {
    //this.adHocTest();

    this.cidViewSelected = action.onViewSelected(this.handleViewSelected);
    this.cidTableSelected = action.onTableSelected(this.handleTableSelected);
    this.cidExecuteQuery = action.onExecuteQuery(this.handleExecuteQuery);
    this.cidExplainQuery = action.onExplainQuery(this.handleExplainQuery);
    this.cidDisconnectDatabase = action.onDisconnectDatabase(this.handleDisconnectDatabase);
    this.cidAlertBar = action.onAlertBar(this.handleAlertBar);
    this.cidSpinner = action.onSpinner(this.handleToggleSpinner);

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
  }

  componentWillUnmount() {
    action.offViewSelected(this.cidViewSelected);
    action.offTableSelected(this.cidTableSelected);
    action.offExecuteQuery(this.cidExecuteQuery);
    action.offExplainQuery(this.cidExplainQuery);
    action.offDisconnectDatabase(this.cidDisconnectDatabase);
    action.offAlertBar(this.cidAlertBar);
    action.offSpinner(this.cidSpinner);
  }

  renderSpinner() {
    if (this.state.spinnerVisible) {
      var spinnerStyle = {
        position: 'fixed',
        top: '50%',
        left: '50%',
        zIndex: '5',
      };

      return <Spinner spinnerName='circle' style={spinnerStyle} />;
    } else {
      return null;
    }
  }

  render() {
    if (!this.state.connected) {
      return (
        <div>
          <div onClick={this.handleCloseAlertBar} >
            { this.state.errorVisible ? <AlertBar errorMessage={this.state.errorMessage}/> : null }
          </div>
          {this.renderSpinner()}
          <ConnectionWindow onDidConnect={this.handleDidConnect} />
        </div>
      );
    }

    var divStyle = { left: this.state.dragBarPosition + 'px' }

    return (
      <div>
        <div onClick={this.handleCloseAlertBar} >
          { this.state.errorVisible ? <AlertBar errorMessage={this.state.errorMessage}/> : null }
        </div>

        {this.renderSpinner()}

        <div>
          <Sidebar
            connectionId={this.state.connectionId}
            tables={this.state.tables}
            selectedTable={this.state.selectedTable}
            selectedTableInfo={this.state.selectedTableInfo}
            databaseName={this.state.databaseName}
            dragBarPosition={this.state.dragBarPosition} />

          <div id="side-dragbar"
            style={divStyle}
            onMouseDown={this.onMouseDown}
            onMouseMove={this.onMouseMove}
            onMouseUp={this.onMouseUp}>
          </div>

          <MainContainer
            results={this.state.results}
            dragBarPosition={this.state.dragBarPosition}
            selectedView={this.state.selectedView} />
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
