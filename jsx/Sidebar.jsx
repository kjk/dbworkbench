/* jshint -W097,-W117 */
'use strict';

var React = require('react');
var ReactDOM = require('react-dom');
var Modal = require('react-modal');

var action = require('./action.js');
var api = require('./api.js');
var Output = require('./Output.jsx');
var filesize = require('filesize');

class Dropdown extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.closeModal = this.closeModal.bind(this);
    this.handleConnection = this.handleConnection.bind(this);
    this.handleActivity = this.handleActivity.bind(this);
    this.handleDisconnect = this.handleDisconnect.bind(this);

    this.state = {
      modalIsOpen: false,
      selectedView: "",
      results: null
    };
  }

  closeModal() {
    this.setState({modalIsOpen: false});
  }

  handleModalCloseRequest() {
    // opportunity to validate something and keep the modal open even if it
    // requested to be closed
    this.closeModal();
  }

  handleConnection() {
    console.log("handleConnection");

    var self = this;
    var connId = this.props.connectionId;
    api.getConnectionInfo(connId, function(data) {
      self.setState({
        results: data,
        modalIsOpen: true,
        selectedView: "Connection Info"
      });
    });
  }

  handleActivity() {
    console.log("handleActivity");

    var self = this;
    var connId = this.props.connectionId;
    api.getActivity(connId, function(data) {
      console.log("getActivity: ", data);
      self.setState({
        results: data,
        modalIsOpen: true,
        selectedView: "Activity"
      });
    });

  }

  handleDisconnect() {
    console.log("handleDisconnect");
    action.disconnectDatabase();
  }

  render() {
    var customStyles = {
      content : {
        display               : 'block',
        overflow              : 'auto',
        top                   : '40%',
        left                  : '50%',
        maxWidth              : '60%',
        transform             : 'translate(-50%, -50%)',
        bottom                : 'none',
        zIndex                : '1'
      }
    };

    var outputStyles = {
      display     :'table-row',
      position    :'absolute',
      padding     :'0',
      margin      :'0',
      top         :'60px',
    };


    var appElement = document.getElementById('main');
    Modal.setAppElement(appElement);


    return (
        <div id="deneme" className='dropdown-window'>
          <div className="list-group">
            <a href="#" className="list-group-item" onClick={this.props.handleRefresh}>Refresh Tables</a>
            <a href="#" className="list-group-item" onClick={this.handleConnection}>Connection Info</a>
            <a href="#" className="list-group-item" onClick={this.handleActivity}>Activity</a>
            <a href="#" className="list-group-item" onClick={this.handleDisconnect}>Disconnect</a>
          </div>


          <Modal
            id='nav'
            isOpen={this.state.modalIsOpen}
            onRequestClose={this.closeModal}
            style={customStyles} >

            <div >
              <div className="modal-header">
                <button type="button" className="close" onClick={this.handleModalCloseRequest}>
                  <span aria-hidden="true">&times;</span>
                  <span className="sr-only">Close</span>
                </button>
                <h4 className="modal-title">{this.state.selectedView}</h4>
              </div>
              <div className="modal-body">

              <Output
                style={outputStyles}
                selectedView={this.state.selectedView}
                results={this.state.results} />
              </div>
            </div>

          </Modal>
        </div>
    );
  }
}

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
    if (!info || $.isEmptyObject(info)) {
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

class Sidebar extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.state = {
      dragging: false,
      tables: [],
    };
  }

  componentWillMount() {
    this.refreshTables();
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

    var self = this;
    api.getTables(connectionId, function(data) {
      // console.log("Refreshing.. " + JSON.stringify(data));
      self.setState({
        tables: data,
      });
    });
  }

  renderTables(tables) {
    var self = this;

    var res = tables.map(function(table) {
      var cls = (table == self.props.selectedTable) ? ' selected' : '';
      var handler = function(e) {
        self.handleSelectTable(e, table);
      };
      return (
        <li onClick={handler} key={table} className={cls}>
          <span><i className='fa fa-table'></i>{table}</span>
        </li>
      );
    });
    return res;
  }

  render() {
    var tables = this.state.tables ? this.renderTables(this.state.tables) : null;
    var divStyle = {
        width: this.props.dragBarPosition + 'px',
    };

    if (this.props.selectedTableInfo != null) {
      var sortList = {
        height: 'calc(100% - 135px)',
      };
    }

    return (
      <div id="sidebar" style={divStyle}>
        <div className="tables-list">
          <div className="wrap">
            <div className="title">
              <i className="fa fa-database"></i>
              <span className="current-database" id="current">{this.props.databaseName}</span>
              <div className='dropdown-menu'>
                <div className="dropdown-cursor">
                  <i className="fa fa-angle-down fa-lg pull-right"></i>
                </div>
                <Dropdown
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

module.exports = Sidebar;
