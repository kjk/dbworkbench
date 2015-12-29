import React from 'react';
import SpinnerCircle from './SpinnerCircle.jsx';
import api from './api.js';
import action from './action.js';

const initialConnectionName = "New connection";

// must match bookmarks.go
const dbTypePostgres = "postgres";
const dbTypeMysql = "mysql";

const defaultPortPostgres = "5432";
const defaultPortMysql = "3306";

const maxBookmarks = 10;

// we need unique ids for unsaved bookmarks. We use negative numbers
// to make sure they don't clash with saved bookmarks (those have positive numbers)
let emptyBookmarkId = -1;

// connecting is async process which might be cancelled
// we use this to uniquely identify connection attempt so that
// whe api.connect() finishes, we can tell if it has been cancelled
// Note: could be state on CoonectionWindow, but we only have one
// of those at any given time so global is just as good
let currConnectionId = 1;

// http://stackoverflow.com/questions/26187189/in-react-js-is-there-any-way-to-disable-all-children-events
// var sayHi = guard("enabled", function(){ alert("hi"); });
// guard.deactivate("enabled");
// sayHi(); // nothing happens
// guard.activate("enabled");
// sayHi(); // shows the alert
var guard = function(key, fn){
  return function(){
    if (guard.flags[key]) {
      return fn.apply(this, arguments);
    }
  };
};

guard.flags = {};
guard.activate = function(key){ guard.flags[key] = true; };
guard.deactivate = function(key){ guard.flags[key] = false; };

function newEmptyBookmark() {
  emptyBookmarkId -= 1;
  // Maybe: change to a class?
  return {
      id: emptyBookmarkId,
      type: dbTypePostgres,
      nick: initialConnectionName,
      database: "",
      url: "",
      host: "",
      user: "",
      password: "",
      port: "" ,
    };
}

class ConnectionWindow extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.newConnectionInfo = this.newConnectionInfo.bind(this);
    this.deleteBookmark = this.deleteBookmark.bind(this);
    this.handleCancel = this.handleCancel.bind(this);
    this.handleConnect = this.handleConnect.bind(this);
    this.handleFormChanged = this.handleFormChanged.bind(this);
    this.handleRememberChange = this.handleRememberChange.bind(this);
    this.selectBookmark = this.selectBookmark.bind(this);
    this.getSelectedBookmark = this.getSelectedBookmark.bind(this);

    // create default bookmark if no bookmarks saved in the backend
    var bookmarks = [newEmptyBookmark()];
    if (gBookmarkInfo && gBookmarkInfo.length > 0) {
      // need to make a copy of the array or else changing bookmark
      // will change gBookmarkInfo
      bookmarks = Array.from(gBookmarkInfo);
    }

    this.state = {
      remember: true,

      connectionErrorMessage: "",
      isConnecting: false,

      bookmarks: bookmarks,
      selectedBookmarkIdx: 0,
    };
  }

  componentDidMount() {
    this.getBookmarks();
  }

  newConnectionInfo(e) {
    var bookmarks = this.state.bookmarks;
    if (bookmarks.length >= maxBookmarks) {
      action.alertBar("Reached connections limit of " + maxBookmarks);
      return;
    }

    bookmarks.push(newEmptyBookmark());
    this.setState({
      bookmarks: bookmarks,
      selectedBookmarkIdx: bookmarks.length-1,
    });
  }

  getSelectedBookmark() {
    return this.state.bookmarks[this.state.selectedBookmarkIdx];
  }

  getBookmarks() {
    api.getBookmarks((data) => {
      console.log("getBookmarks: ", data);
      if (!data) {
          data = [newEmptyBookmark()];
      }
      this.setState({
        bookmarks: data,
      });
    });
  }

  deleteBookmark(e) {
    e.stopPropagation();

    const idStr = e.target.attributes["data-custom-attribute"].value;
    const id = parseInt(idStr, 10);
    // bookmarks with negative id are not yet saved (only exist in the frontend)
    if (id < 0) {
      let selectedIdx = this.state.selectedBookmarkIdx;
      let bookmarks = this.state.bookmarks.filter((b) => b.id != id);
      if (selectedIdx >= bookmarks.length) {
        selectedIdx = bookmarks.length - 1;
      }
      if (bookmarks.length == 0) {
        bookmarks = [newEmptyBookmark()];
        selectedIdx = 0;
      }
      this.setState({
        bookmarks: bookmarks,
        selectedBookmarkIdx: selectedIdx,
      });
      return;
    }

    api.removeBookmark(id, (data) => {
      console.log("deleteBookmarks removing: ", id, " data: ", data);

      let bookmarks = [newEmptyBookmark()];
      let selectedBookmarkIdx = 0;
      if (data !== undefined && data.length > 0) {
        bookmarks = data;
      }

      this.setState({
        bookmarks: bookmarks,
        selectedBookmarkIdx: selectedBookmarkIdx,
      });
    });
  }

  selectBookmark(e) {
    e.stopPropagation();

    var idxStr = e.currentTarget.attributes["data-custom-attribute"].value;
    var idx = parseInt(idxStr, 10);
    console.log("selectBookmark, idx:", idx);

    this.setState({
      selectedBookmarkIdx: idx,
      connectionErrorMessage: ""
    });
  }

  handleFormChanged(name, e) {
    //console.log("handleFormChanged: name=", name, " val=", e.target.value);

    let b = this.getSelectedBookmark();
    let prevDatabase = b["database"];
    b[name] = e.target.value;
    let dbName = b["database"];

    // if nick has not been modified by user, make it equal to database name
    let nick = b["nick"];
    if ((nick == initialConnectionName) || (nick == prevDatabase)) {
      if (dbName != "") {
        b["nick"] = dbName;
      }
    }

    const bookmarks = this.state.bookmarks;
    bookmarks[this.selectedBookmarkIdx] = b;
    this.setState({
      bookmarks: bookmarks,
      connectionErrorMessage: "",
    });
  }

  handleRememberChange(e) {
    var newRemeber = !this.state.remember;
    this.setState({
      remember: newRemeber,
    });
    //console.log("remember changed to: " + newRemeber);
  }

  handleConnect(e) {
    e.preventDefault();
    console.log("handleConnect");

    let b = this.getSelectedBookmark();

    let id = b["id"];
    let nick = b["nick"];
    let dbType = b["type"];
    let host = b["host"];
    let port = b["port"];
    let user = b["user"];
    let pass = b["password"];
    let db = b["database"];
    let rememberConnection = this.state.remember;

    let url = "";
    let urlSafe = "";
    if (dbType == dbTypePostgres) {
      if (port.length == 0) {
        port = defaultPortPostgres;
      }
      url = "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + db;
      urlSafe = url;
      if (pass != "") {
        urlSafe = "postgres://" + user + ":" + "***" + "@" + host + ":" + port + "/" + db;
      }
    }
    else if (dbType == dbTypeMysql)
    {
      // mysql format:
      // username:password@protocol(address)/dbname?param=value
      // dbname can be empty
      if (port.length == 0) {
        port = defaultPortMysql;
      }
      // pareTime: conver time from []byte to time.Time
      // https://github.com/go-sql-driver/mysql#parsetime
      url = user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + db + "?parseTime=true";
      urlSafe = url;
      if (pass != "") {
        urlSafe = user + ":" + "***" + "@tcp(" + host + ":" + port + ")/" + db + "?parseTime=true";
      }
    } else {
      console.log("invalid type: " + dbType);
      // TODO: how to error out?
    }

    console.log("URL:" + url);
    this.setState({
      isConnecting: true,
    });
    const myConnectionId = currConnectionId;
    api.connect(dbType, url, urlSafe, (resp) => {
      if (myConnectionId != currConnectionId) {
        console.log("ignoring completion of a cancelled connection");
        return;
      }
      ++currConnectionId;
      if (resp.error) {
        console.log("handleConnect: resp.error: ", resp.error);

        this.setState({
          connectionErrorMessage: resp.error,
          isConnecting: false,
        });
        return;
      }

      b = this.getSelectedBookmark();
      if (!rememberConnection) {
        console.log("did connect, not saving a bookmark");
        return;
      }
      console.log("did connect, saving a bookmark " + b);
      api.addBookmark(b, (data) => {
        const connId = resp.ConnectionID;
        const connStr = url;
        const databaseName = resp.CurrentDatabase;
        const capabilities = resp.Capabilities;
        this.props.onDidConnect(connStr, connId, databaseName, capabilities);
      });
    });
  }

  handleCancel(e) {
    e.preventDefault();
    console.log("handleCancel");
    // to tell api.connec() callback that we've been cancelled
    ++currConnectionId;

    this.setState({
      isConnecting: false
    });
  }

  renderErrorOptional(errorText) {
    if (errorText != "") {
      return <div className="col-md-12 connection-error">Error: {errorText}</div>;
    }
  }

  renderBookMarks() {
    guard.activate("bookmarksEnabled");
    if (this.state.isConnecting) {
      guard.deactivate("bookmarksEnabled");
    }

    let bookmarks = [];
    for (var i = 0; i < this.state.bookmarks.length; i++) {
      let b = this.state.bookmarks[i];
      let id = b["id"];
      let nick = b["nick"];

      let className = "list-group-item";
      if (i == this.state.selectedBookmarkIdx) {
        className = "list-group-item active";
      }

      bookmarks.push(
        <a key={id} data-custom-attribute={i} href="#" className={className} onClick={guard("bookmarksEnabled", this.selectBookmark)}>
          {nick}
          <i data-custom-attribute={id} onClick={guard("bookmarksEnabled", this.deleteBookmark)} className="fa fa-times pull-right"></i>
        </a>
      );
    }

    return (
      <div className="list-group list-special">
        <a href="#" className="list-group-item title" onClick={guard("bookmarksEnabled", this.newConnectionInfo)} >
          Connections
          <i className="fa fa-plus pull-right"></i>
        </a>

        <hr/>

        {bookmarks}
      </div>
    );
  }

  renderFormElements() {
    let b = this.getSelectedBookmark();

    let dbType = b["type"];
    let defaultPort = "0";
    if (dbType == dbTypePostgres) {
      defaultPort = defaultPortPostgres;
    } else if (dbType == dbTypeMysql) {
      defaultPort = defaultPortMysql;
    } else {
      console.log("Unknown type: " + dbType);
    }

    let disable = this.state.isConnecting;

    return (
      <div>
        <div className="col-md-8">
          <div className="form-group">
            <label className="control-label" htmlFor="db_nickname">Nickname</label>
            <input
              type="text"
              id="db_nickname"
              className="form-control input-sm"
              value = {b["nick"]}
              disabled={disable}
              onChange={this.handleFormChanged.bind(this, 'nick')}/>
          </div>
        </div>

        <div className="col-md-4">
          <div className="form-group">
            <label className="control-label" htmlFor="db_type">Type</label>
            <select
              id="db_type"
              className="form-control input-sm"
              value={dbType}
              disabled={disable}
              onChange={this.handleFormChanged.bind(this, 'type')}>
                <option value={dbTypePostgres}>PostgreSQL</option>
                <option value={dbTypeMysql}>MySQL</option>
            </select>
          </div>
        </div>

        <div className="col-md-8">
          <div className="form-group">
            <label className="control-label" htmlFor="db_hostname">Hostname</label>
            <input
              type="text"
              id="db_hostname"
              className="form-control input-sm"
              value = {b["host"]}
              disabled={disable}
              onChange={this.handleFormChanged.bind(this, 'host')}/>
          </div>
        </div>

        <div className="col-md-4">
          <div className="form-group">
            <label className="control-label" htmlFor="db_port">Port</label>
            <input
              type="text"
              id="db_port"
              className="form-control input-sm"
              value = {b["port"]}
              disabled={disable}
              onChange={this.handleFormChanged.bind(this, 'port')} placeholder={defaultPort}/>
          </div>
        </div>

        <div className="col-md-12">
          <div className="form-group">
            <label className="control-label" htmlFor="db_database">Database</label>
            <input
              type="text"
              id="db_database"
              className="form-control input-sm"
              value = {b["database"]}
              disabled={disable}
              onChange={this.handleFormChanged.bind(this, 'database')} />
          </div>
        </div>

        <div className="col-md-6">
          <div className="form-group">
            <label className="control-label" htmlFor="db_user">User</label>
            <input
              type="text"
              id="db_user"
              className="form-control input-sm"
              value = {b["user"]}
              disabled={disable}
              onChange={this.handleFormChanged.bind(this, 'user')} />
          </div>
        </div>

        <div className="col-md-6">
          <div className="form-group">
            <label className="control-label" htmlFor="db_pass">Password</label>
            <input
              type="password"
              id="db_pass"
              className="form-control input-sm"
              value = {b["password"]}
              disabled={disable}
              onChange={this.handleFormChanged.bind(this, 'password')} />
          </div>
        </div>

        <div className="col-md-12 right">
          <label className="control-label" htmlFor="pwd-remember">
            <input type="checkbox"
              id="pwd-remember"
              checked={this.state.remember}
              disabled={disable}
              onChange={this.handleRememberChange}
            /> Remember
         </label>
        </div>

        <div className="col-md-12 right light-text smaller-text">
          <i className="fa fa-lock fa1"></i>&nbsp;Database crendentials are stored
          securely on your computer
        </div>

        <div className="col-md-12">
          &nbsp;&nbsp;
        </div>

        {this.renderErrorOptional(this.state.connectionErrorMessage)}

        {this.renderConnectOrCancel()}
      </div>
    );
  }

  renderConnectOrCancel() {
    let styleDiv = {
      position: 'relative'
    };

    let styleSpinner = {
      zIndex: '5',
      position: 'absolute',
      right: '-32px',
      top: '8'
    };


    if (this.state.isConnecting) {
      return (
        <div className="col-md-12" style={styleDiv}>
          <button onClick={this.handleCancel} className="btn btn-block btn-danger small">Cancel</button>
          <SpinnerCircle visible={true} style={styleSpinner}/>
        </div>
      );
    }

    return (
      <div className="col-md-12">
        <button onClick={this.handleConnect} className="btn btn-block btn-primary small">Connect</button>
      </div>
    );
  }

  renderForm() {
    if (this.state.selectedBookmarkIdx >= 0) {
      return (
        <form role="form">
          {this.renderFormElements()}
        </form>
      );
    }

    // TODO: I don't think it ever happens
    var imageStyle = {
      width: "30%",
      height: "30%"
    };

    return (
      <form role="form">
        <div className="col-md-12 text-center">
            <img class="img-responsive center-block small"
              src="/s/img/icon.png"
              alt="" style={imageStyle}/>
            <h5>Please add a connection</h5>
        </div>
      </form>
    );
  }

  renderConnectionWindow() {
    return (
      <div className="container small">
        <div className="row">

          <div className="col-md-4">
            {this.renderBookMarks()}
          </div>

          <div className="col-md-8">
            {this.renderForm()}
          </div>

        </div>
      </div>
    );
  }

  render() {
    let versionStyle = {
      position: 'absolute',
      bottom: '0px',
      right: '0',
      padding: '5px',
      fontSize: '12px',
      color: '#A9A9A9',
    };

    return (
      <div id="connection-window">
          <div className='logo-container'><img className='resize_fit_center' src='/s/img/dbhero-sm.png' /></div>
          <div className="connection-settings">
            {this.renderConnectionWindow()}
            <hr/>
          </div>
          <div style={versionStyle}>Version: {gVersionNumber}</div>
      </div>
    );
  }
}

module.exports = ConnectionWindow;
