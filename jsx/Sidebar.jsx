/* jshint -W097,-W117 */
'use strict';

var api = require('./api.js');

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

module.exports = Sidebar;
