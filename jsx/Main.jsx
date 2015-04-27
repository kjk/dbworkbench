/* jshint -W097,-W117 */
'use strict';

var api = require('./api.js');

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

var TableInformation = React.createClass({
  renderTableInfo: function(info) {
    console.log("renderTableInfo: info: ", info);
    if (info && !$.isEmptyObject(info)) {
      return (
        <ul>
          <li>Size: <span>{info.total_size}</span></li>
          <li>Data size: <span>{info.data_size}</span></li>
          <li>Index size: <span>{info.index_size}</span></li>
          <li>Estimated rows: <span>{info.rows_count}</span></li>
        </ul>
      );
    }
  },

  render: function() {
    var info = this.renderTableInfo(this.props.tableInfo);
    return (
      <div className="table-information">
        <div className="wrap">
          <div className="title">Table Information</div>
          {info}
        </div>
      </div>
    );
  }
});

// TODO: pass connectionId and use it in api calls
var Sidebar = React.createClass({

  getInitialState: function() {
    return {
      tableNames: [],
      selectedTableName: "",
      selectedTableInfo: null,
    };
  },

  componentDidMount: function() {
    var self = this;
    api.getTables(function(data) {
      //console.log("componentDidMount: ", data);
      self.setState({
        tableNames: data,
        tableInfo: null,
      });
    });
  },

  handleSelectTable: function(e) {
    e.preventDefault();
    var table = e.target.textContent.trim();
    console.log("handleSelectTable: ", e.target, " table:", table);
    var self = this;
    api.call("get", "/tables/" + table + "/info", {}, function(data) {
      console.log("handleSelectTable: tableInfo: ", data);
      self.setState({
        selectedTableInfo: data,
        selectedTableName: table,
      });
    });

    var sortColumn = null;
    var sortOrder = null;
    var params = { limit: 100, sort_column: sortColumn, sort_order: sortOrder };
    api.getTableRows(table, params, function(data) {
      console.log("handleSelectTable: got table rows: ", data);
      self.props.onGotResults(data);
    });
  },

  // TODO: remove id="tabels"
  render: function() {
    var self = this;
    var tables = this.state.tableNames.map(function(item) {
      var cls;
      if (item == self.state.selectedTableName) {
        cls += ' selected';
      }
      return (
        <li onClick={self.handleSelectTable} key={item} className={cls}>
          <span><i className='fa fa-table'></i>{item}</span>
        </li>
      );
    });

    var tableInfo = this.state.selectedTableInfo;
    return (
      <div id="sidebar">
        <div className="tables-list">
          <div className="wrap">
            <div className="title">
              <i className="fa fa-database"></i> <span className="current-database" id="current_database"></span>
              <span className="refresh" id="refresh_tables" title="Refresh tables list"><i className="fa fa-refresh"></i></span>
            </div>
            <ul id="tables">
              {tables}
            </ul>
          </div>
        </div>
        <TableInformation tableInfo={tableInfo} />
      </div>
    );
  }

});

var Input = React.createClass({
  render: function() {
    return (
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
    );
  }
});

/*
function buildTable(results, sortColumn, sortOrder) {
  results.columns.forEach(function(col) {
      if (col === sortColumn) {
          cols += "<th data='" + col + "'" + "data-sort-order=" + sortOrder + ">" + col + "&nbsp;" + sortArrow(sortOrder) + "</th>";
      } else {
          cols += "<th data='" + col + "'>" + col + "</th>";
      }
  });
}
*/

var Output = React.createClass({
  renderCols: function(columns, sortColumn, sortOrder) {
    columns.map(function(col) {
      // TODO: use sortColumn and sortOrder
      return (
        <th data={col}>{col}</th>
      );
    });
  },

  renderRows: function(rows) {
    var children = rows.map(function(row) {
      return (
        <td><div>{row}</div></td>
      );
    });

    return (
      <tr>{children}</tr>
    );
  },

  renderResults: function(results) {
    var cols = this.renderCols(results.columns);
    var rows = this.renderRows(results.rows);
    return (
      <table id="results" className="table">
        <thead>{cols}</thead>
        <tbody>{rows}</tbody>
      </table>
    );
  },

  renderNoResults: function() {
    return (
      <table id="results" className="table">
        <tr><td>No records found</td></tr>
      </table>
    );
  },

  renderError: function(errorMsg) {
    return (
      <table id="results" className="table">
        <tr><td>ERROR: {errorMsg}</td></tr>
      </table>
    );
  },

  render: function() {
    var cls = "table";
    var results;
    if (!this.props.results) {
      cls += " empty";
      results = this.renderNoResults();
    } else {
      if (this.props.results.error) {
        cls += " empty";
        results = this.renderError(this.props.results.error);
      } else {
        results = this.renderResults(this.props.results);
      }
    }

    return (
      <div id="output">
        <div className="wrapper">
            {results}
        </div>
      </div>
    );
  }
});

module.exports.TopNav = TopNav;
module.exports.Sidebar = Sidebar;
module.exports.Input = Input;
module.exports.Output = Output;
