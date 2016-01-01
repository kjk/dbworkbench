import React from 'react';
import ReactDOM from 'react-dom';
import DatabaseMenuDropdown from './DatabaseMenuDropdown.jsx'
import filesize from 'filesize';
import * as action from './action.js';
import * as api from './api.js';
import * as store from './store.js';

function isEmptyObject(object) {
  let name;
  for (name in object) {}
  return name === undefined;
};

class TableInformation extends React.Component {
  renderTableInfo(info) {
    const dataSize = parseInt(info.data_size);
    const dataSizePretty = filesize(dataSize);
    const indexSize = parseInt(info.index_size);
    const indexSizePretty = filesize(indexSize);
    const totalSize = dataSize + indexSize;
    const totalSizePretty = filesize(totalSize);
    const rowCount = parseInt(info.rows_count);

    // TODO: better done as a class,maybe on parent element
    const style = {
      backgroundColor: "white",
    };

    return (
      <ul style={style}>
        <li><span className="table-info-light">Size: </span><span>{totalSizePretty}</span></li>
        <li><span className="table-info-light">Data size: </span><span>{dataSizePretty}</span></li>
        <li><span className="table-info-light">Index size: </span><span>{indexSizePretty}</span></li>
        <li><span className="table-info-light">Estimated rows: </span><span>{rowCount}</span></li>
      </ul>
    );
  }

  renderTableInfoContainer() {
    const info = this.props.tableInfo;
    if (!info || isEmptyObject(info)) {
      return;
    }

    const tableInfo = this.renderTableInfo(info);
    return (
      <div className="wrap">
          <div className="title">
          <i className="fa fa-info"></i>
          <span className="current-table-information">Table Information</span></div>
          {tableInfo}
      </div>
    );
  }

  render() {
    return (
      <div className="table-information">
        {this.renderTableInfoContainer()}
      </div>
    );
  }
}

export default class Sidebar extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.sidebarDx = store.getSidebarDx();
    this.state = {
      tables: [],
    };
  }

  componentWillMount() {
    this.refreshTables();
    this.cidSidebarDx = store.onSidebarDx( (dx) => {
      this.sidebarDx = dx;
      const el = ReactDOM.findDOMNode(this);
      el.style.width = dx + "px";
    });
  }

  componentWillUnmount() {
    store.offSidebarDx(this.cidSidebarDx);
  }

  handleRefreshDatabase(e) {
    e.preventDefault();

    console.log("handleRefreshDatabase");

    // TODO: make some kind of UI representation of refresh
    // just to show users that the action was successful.
    this.refreshTables();;
  }

  handleSelectTable(e, table) {
    e.preventDefault();
    action.tableSelected(table);
  }

  refreshTables() {
    var connectionId = this.props.connectionId;

    api.getTables(connectionId, (data) => {
      // console.log("Refreshing.. " + JSON.stringify(data));
      this.setState({
        tables: data,
      });
    });
  }

  renderTables(tables) {
    if (!tables) {
      return null;
    }
    const selectedTable = this.props.selectedTable;
    const res = tables.map((table) => {
      const cls = (table == selectedTable) ? ' selected' : '';
      let handler = (e) => this.handleSelectTable(e, table);
      return (
        <li onClick={handler} key={table} className={cls}>
          <span><i className='fa fa-table'></i>{table}</span>
        </li>
      );
    });
    return res;
  }

  render() {
    // TODO: on database connect gets rendered 28 times
    //console.log("Sidebar render");

    var tables = this.renderTables(this.state.tables);
    var style = {
        width: this.sidebarDx,
    };

    if (this.props.selectedTableInfo != null) {
      var sortList = {
        height: 'calc(100% - 135px)',
      };
    }

    return (
      <div id="sidebar" style={style}>
        <div className="tables-list">
          <div className="wrap">
            <div className="title">
              <i className="fa fa-database"></i>
              <span className="current-database" id="current">{this.props.databaseName}</span>
              <div className='dropdown-menu'>
                <div className="dropdown-cursor">
                  <i className="fa fa-angle-down fa-lg pull-right"></i>
                </div>
                <DatabaseMenuDropdown
                  connectionId={this.props.connectionId}
                  handleRefresh={this.handleRefreshDatabase.bind(this)} />
              </div>
            </div>
            <ul style={sortList}>
              {tables}
            </ul>
          </div>
        </div>
        <TableInformation tableInfo={this.props.selectedTableInfo} />


      </div>
    );
  }
}

