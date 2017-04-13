import React from "react";
import ReactDOM from "react-dom";
import PropTypes from "prop-types";
import DatabaseMenuDropdown from "./DatabaseMenuDropdown.jsx";
import TableInformation from "./TableInformation.jsx";
import * as action from "./action.js";
import * as api from "./api.js";
import * as store from "./store.js";

export default class Sidebar extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.handleRefreshTables = this.handleRefreshTables.bind(this);

    this.sidebarDx = store.getSidebarDx();
    this.state = {
      tables: [],
    };
  }

  componentWillMount() {
    this.refreshTables();
    store.onSidebarDx(dx => {
      this.sidebarDx = dx;
      const el = ReactDOM.findDOMNode(this);
      el.style.width = dx + "px";
    }, this);
  }

  componentWillUnmount() {
    store.offAllForOwner(this);
  }

  handleRefreshTables(e) {
    //console.log('handleRefreshTables');
    //e.preventDefault();
    // TODO: make some kind of UI representation of refresh
    // just to show users that the action was successful.
    this.refreshTables();
  }

  handleSelectTable(e, table) {
    e.preventDefault();
    action.tableSelected(table);
  }

  refreshTables() {
    this.props.refreshAllTableInformation();
  }

  renderTables(tables) {
    if (!tables) {
      return null;
    }
    const selectedTable = this.props.selectedTable;
    const res = tables.map(table => {
      const cls = table == selectedTable ? " selected" : "";
      let handler = e => this.handleSelectTable(e, table);
      return (
        <li onClick={handler} key={table} className={cls}>
          <span><i className="fa fa-table" />{table}</span>
        </li>
      );
    });
    return res;
  }

  render() {
    // TODO: on database connect gets rendered 28 times
    //console.log("Sidebar render");

    const tables = this.renderTables(this.props.tables);
    const style = {
      width: this.sidebarDx,
    };

    let sortList = {};
    const tableInfo = this.props.selectedTableInfo;
    if (tableInfo != null) {
      sortList = {
        height: "calc(100% - 135px)",
      };
    }

    return (
      <div id="sidebar" style={style}>
        <div className="tables-list">
          <div className="wrap">
            <div className="title">
              <i className="fa fa-database" />
              <span className="current-database" id="current">
                {this.props.databaseName}
              </span>
              <div className="dropdown-menu">
                <div className="dropdown-cursor">
                  <i className="fa fa-angle-down fa-lg pull-right" />
                </div>
                <DatabaseMenuDropdown
                  connectionId={this.props.connectionId}
                  onRefreshTables={this.handleRefreshTables}
                />
              </div>
            </div>
            <ul style={sortList}>
              {tables}
            </ul>
          </div>
        </div>
        {tableInfo
          ? <TableInformation tableInfo={this.props.selectedTableInfo} />
          : null}
      </div>
    );
  }
}

Sidebar.propTypes = {
  refreshAllTableInformation: PropTypes.func,
  selectedTable: PropTypes.string,
  tables: PropTypes.array, // TODO: more specific
  selectedTableInfo: PropTypes.any, // TODO: more specific
  databaseName: PropTypes.string,
  connectionId: PropTypes.number,
};
