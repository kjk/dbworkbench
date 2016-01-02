import React from 'react';
import ReactDOM from 'react-dom';
import DatabaseMenuDropdown from './DatabaseMenuDropdown.jsx';
import TableInformation from './TableInformation.jsx';
import filesize from 'filesize';
import * as action from './action.js';
import * as api from './api.js';
import * as store from './store.js';

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
    store.onSidebarDx((dx) => {
      this.sidebarDx = dx;
      const el = ReactDOM.findDOMNode(this);
      el.style.width = dx + 'px';
    }, this);
  }

  componentWillUnmount() {
    store.offAllForOwner(this);
  }

  handleRefreshDatabase(e) {
    e.preventDefault();

    console.log('handleRefreshDatabase');

    // TODO: make some kind of UI representation of refresh
    // just to show users that the action was successful.
    this.refreshTables();;
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
    const res = tables.map((table) => {
      const cls = (table == selectedTable) ? ' selected' : '';
      let handler = (e) => this.handleSelectTable(e, table);
      return (
        <li onClick={ handler } key={ table } className={ cls }>
          <span><i className='fa fa-table'></i>{ table }</span>
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
    if (this.props.selectedTableInfo != null) {
      sortList = {
        height: 'calc(100% - 135px)',
      };
    }

    return (
      <div id="sidebar" style={ style }>
        <div className="tables-list">
          <div className="wrap">
            <div className="title">
              <i className="fa fa-database"></i>
              <span className="current-database" id="current">{ this.props.databaseName }</span>
              <div className='dropdown-menu'>
                <div className="dropdown-cursor">
                  <i className="fa fa-angle-down fa-lg pull-right"></i>
                </div>
                <DatabaseMenuDropdown connectionId={ this.props.connectionId } handleRefresh={ this.handleRefreshDatabase.bind(this) } />
              </div>
            </div>
            <ul style={ sortList }>
              { tables }
            </ul>
          </div>
        </div>
        <TableInformation tableInfo={ this.props.selectedTableInfo } />
      </div>
      );
  }
}

