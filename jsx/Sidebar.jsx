/* jshint -W097,-W117 */
'use strict';

var React = require('react');
var Modal = require('react-modal');

var action = require('./action.js');
var api = require('./api.js');

class Dropdown extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.closeModal = this.closeModal.bind(this);
    this.handleConnection = this.handleConnection.bind(this);
    this.handleDisconnect = this.handleDisconnect.bind(this);
    this.openModal = this.openModal.bind(this);
    this.state = { modalIsOpen: false };
  }

  openModal() {
    this.setState({modalIsOpen: true});
  }

  closeModal() {
    this.setState({modalIsOpen: false});
  }

  handleConnection() {
    console.log("handleConnection");
  }

  handleActivity() {
    console.log("handleActivity");
  }

  handleDisconnect() {
    console.log("handleDisconnect");
    action.disconnectDatabase();
  }

  render() {
    const customStyles = {
      content : {
        top                   : '50%',
        left                  : '50%',
        right                 : 'auto',
        bottom                : 'auto',
        marginRight           : '-50%',
        transform             : 'translate(-50%, -50%)'
      }
    };

    Modal.setAppElement('#deneme');

    return (
        <div id="deneme" className='dropdown-window'>
          <div className="list-group">
            <a href="#" className="list-group-item" onClick={this.handleConnection}>Connection</a>
            <a href="#" className="list-group-item" onClick={this.openModal}>Activity</a>
            <a href="#" className="list-group-item" onClick={this.handleDisconnect}>Disconnect</a>
          </div>

          <Modal
            isOpen={this.state.modalIsOpen}
            onRequestClose={this.closeModal}
            style={customStyles} >

            <h2>Hello</h2>
            <button onClick={this.closeModal}>close</button>
            <div>I am a modal</div>
            <form>
              <input />
              <button>tab navigation</button>
              <button>stays</button>
              <button>inside</button>
              <button>the modal</button>
            </form>
          </Modal>

        </div>
    )
  }
}

class TableInformation extends React.Component {
  renderTableInfo(info) {
    if (info && !$.isEmptyObject(info)) {
      return (
        <ul>
          <li><span className="table-info-light">Size: </span><span>{info.total_size}</span></li>
          <li><span className="table-info-light">Data size: </span><span>{info.data_size}</span></li>
          <li><span className="table-info-light">Index size: </span><span>{info.index_size}</span></li>
          <li><span className="table-info-light">Estimated rows: </span><span>{info.rows_count}</span></li>
        </ul>
      );
    }
  }

  renderTableInfoContainer() {
    var info = this.renderTableInfo(this.props.tableInfo);
    if (info) {
      return (
        <div className="wrap">
          <div className="title">
            <i className="fa fa-info"></i>
            <span className="current-table-information">Table Information</span></div>
            {info}
        </div>
      );
    } else {
      return (<div></div>);
    }

  }

  render() {
    var info = this.renderTableInfo(this.props.tableInfo);
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

  // TODO: remove id="tables"
  render() {
    var tables = this.state.tables ? this.renderTables(this.state.tables) : null;
    var divStyle = {
        width: this.props.dragBarPosition + 'px',
    }

    // <span className="refresh" id="refresh_tables"
    //                 title="Refresh tables list" onClick={this.handleRefreshDatabase}> <i className="fa fa-refresh"></i>
    //           </span>


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
                <Dropdown />
              </div>
            </div>
            <ul id="tables">
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
