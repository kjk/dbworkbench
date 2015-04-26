/* jshint -W097,-W117 */
'use strict';

var ViewContent = 0;
var ViewStructure = 1;
var ViewIndexes = 2;
var ViewSQLQuery = 3;
var ViewHistory = 4;
var ViewActivity = 5;
var ViewConnection = 6;

var allViews = [0, 1, 2, 3, 4, 5, 6];

var viewNames = [
  "Content",
  "Structure",
  "Indexes",
  "SQL Query",
  "History",
  "Activity",
  "Connection"
];

var TopNav = React.createClass({
  getInitialState: function() {
    return {
      view: ViewSQLQuery
    };
  },

  switchToView: function(view) {
    this.setState({
      view: view,
    });
  },

  switchToContent: function() { this.switchToView(ViewContent); },
  switchToStructure: function() { this.switchToView(ViewStructure); },
  switchToIndexes: function() { this.switchToView(ViewIndexes); },
  switchToSQLQuery: function() { this.switchToView(ViewSQLQuery); },
  switchToHistory: function() { this.switchToView(ViewHistory); },
  switchToActivity: function() { this.switchToView(ViewActivity); },
  switchToConnection: function() { this.switchToView(ViewConnection); },

  render: function() {
    var handlers = [
      this.switchToContent,
      this.switchToStructure,
      this.switchToIndexes,
      this.switchToSQLQuery,
      this.switchToHistory,
      this.switchToActivity,
      this.switchToConnection
    ];
    var self = this;

    var children = allViews.map(function(view) {
      var cls;
      if (self.state.view == view) {
        cls = "selected";
      }
      var handler = handlers[view];
      var txt = viewNames[view];
      return <li onClick={handler} key={view} className={cls}>{txt}</li>;
    });

    return (
      <div id="nav">
        <ul>
          {children}
        </ul>

        <a href="#" id="edit_connection" className="btn btn-primary btn-sm"><i className="fa fa-gear"></i> Edit Connection</a>
      </div>
    );
  },
});

var Sidebar = React.createClass({

  render: function() {
    return (
      <div id="sidebar">
        <div className="tables-list">
          <div className="wrap">
            <div className="title">
              <i className="fa fa-database"></i> <span className="current-database" id="current_database"></span>
              <span className="refresh" id="refresh_tables" title="Refresh tables list"><i className="fa fa-refresh"></i></span>
            </div>
            <ul id="tables"></ul>
          </div>
        </div>
        <div className="table-information">
          <div className="wrap">
            <div className="title">Table Information</div>
            <ul>
              <li>Size: <span id="table_total_size"></span></li>
              <li>Data size: <span id="table_data_size"></span></li>
              <li>Index size: <span id="table_index_size"></span></li>
              <li>Estimated rows: <span id="table_rows_count"></span></li>
            </ul>
          </div>
        </div>
      </div>
    );
  }

});

var Body = React.createClass({

  render: function() {
    return (
      <div id="body">
        <div id="input">
          <div className="wrapper">
            <div id="custom_query"></div>
            <div className="actions">
              <input type="button" id="run" value="Run Query" className="btn btn-sm btn-primary" />
              <input type="button" id="explain" value="Explain Query" className="btn btn-sm btn-default" />
              <input type="button" id="csv" value="Download CSV" className="btn btn-sm btn-default" />

              <div id="query_progress">Please wait, query is executing...</div>
            </div>
          </div>
        </div>
        <div id="output">
          <div className="wrapper">
            <table id="results" className="table"></table>
          </div>
        </div>
      </div>
    );
  }
});

module.exports.TopNav = TopNav;
module.exports.Sidebar = Sidebar;
module.exports.Body = Body;
