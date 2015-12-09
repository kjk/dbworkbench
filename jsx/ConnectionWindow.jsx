/* jshint -W097,-W117 */
'use strict';

var React = require('react');

var api = require('./api.js');
var action = require('./action.js');
var _ = require('underscore');

var initName = "New connection";

function newEmptyBookmark() {
  return {
      id: 0,
      url: "",
      host: "",
      database: "New connection",
      user: "",
      password: "",
      port: "" ,
      ssl: ""
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
    this.selectBookmark = this.selectBookmark.bind(this);
    this.getActiveBookmark = this.getActiveBookmark.bind(this);

    // if there are no bookmarks defined on the server,
    // add a default bookmark
    var bookmarks = gBookmarkInfo
    if (!bookmarks || bookmarks.Length == 0) {
      bookmarks = [newEmptyBookmark()];
    }
    this.state = {
      connectionErrorMessage: "",
      isConnecting: false,

      bookmarks: bookmarks,
      activeBookmark: 0,
    };
  }

  newConnectionInfo(e) {
    var bookmarkLimit = 10;
    var bookmarks = this.state.bookmarks;
    if (bookmarks.length >= bookmarkLimit) {
      action.alertBar("Reached connections limist of " + bookmarkLimit);
      return;
    }

    bookmarks.push(newEmptyBookmark())
    this.setState({
      bookmarks: bookmarks,
      activeBookmark: bookmarks.length-1,
    });
  }

  getBookmark(e) {
    var self = this;
    api.getBookmarks(function(data) {
      console.log("getBookmarks: ", data);

      self.setState({
        bookmarks: data,
      });
    });
  }

  deleteBookmark(e) {
    e.stopPropagation();
    var self = this;

    var id = e.target.attributes["data-custom-attribute"].value;
    api.removeBookmark(id, function(data) {
      console.log("deleteBookmarks removing: ", id, " data: ", data);

      var bookmarks = [newEmptyBookmark()];
      var activeBookmark = 0;
      if (data !== undefined && data.length > 0) {
        bookmarks = data;
      }

      self.setState({
        bookmarks: bookmarks,
        activeBookmark: activeBookmark,
      });
    });
  }

  selectBookmark(e) {
    e.stopPropagation();

    var idxStr = e.currentTarget.attributes["data-custom-attribute"].value;
    var idx = parseInt(idxStr, 10);
    console.log("selectBookmark", idx);

    this.setState({
      activeBookmark: idx,
      connectionErrorMessage: ""
    });
  }

  handleFormChanged(name, e) {
    console.log("handleFormChanged: ", e.target.value);

    var change = this.state.bookmarks;

    if (change[this.state.activeBookmark]['oldDatabase'] === undefined) {
      change[this.state.activeBookmark]['oldDatabase'] = change[this.state.activeBookmark]['database'];
      console.log("here", change)
    }
    change[this.state.activeBookmark][name] = e.target.value;

    this.setState({
      bookmarks: change,
      connectionErrorMessage: "",
    });
  }

  getActiveBookmark() {
    return this.state.bookmarks[this.state.activeBookmark];
  }

  handleConnect(e) {
    e.preventDefault();
    console.log("handleConnect");

    var b = this.getActiveBookmark();

    var id = b["id"];
    var type = b["type"];
    var host = b["host"];
    var port = b["port"];
    var user = b["user"];
    var pass = b["password"];
    var db = b["database"];
    var ssl = "disable";

    if (port.length == 0) {
      port = "5432";
    }

    var url = "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + db + "?sslmode=" + ssl;

    console.log("URL:" + url);
    var self = this;
    this.setState({
      isConnecting: true,
    });
    api.connect(url, function(resp) {
      if (resp.error) {
        console.log("handleConnect: resp.error: ", resp.error);

        self.setState({
          connectionErrorMessage: resp.error,
          isConnecting: false,
        });
      } else {
        b = self.getActiveBookmark()
        console.log("did connect, saving bookmark " + b);
        api.addBookmark(b);

        var connId = resp.ConnectionID;
        var connStr = url;
        var databaseName = resp.CurrentDatabase;
        self.props.onDidConnect(connStr, connId, databaseName);
      }
    });
  }

  handleCancel(e) {
    e.preventDefault();
    console.log("handleCancel");
  }

  renderError(errorText) {
    return (
      <div className="alert alert-danger">{errorText}</div>
    );
  }

  renderBookMarks() {
    var bookmarks = [];
    for (var i = 0; i < this.state.bookmarks.length; i++) {
      var bookmark = this.state.bookmarks[i];
      var name = bookmark["database"];
      var id = bookmark["id"]

      var removeButton = <i data-custom-attribute={id} onClick={this.deleteBookmark} className="fa fa-times pull-right"></i>;

      var className = "list-group-item"
      if (i == this.state.activeBookmark) {
        className = "list-group-item active"
      }

      bookmarks.push(
        <a key={id} data-custom-attribute={i} href="#" className={className} onClick={this.selectBookmark}>
          {name}
          {removeButton}
        </a>
      );
    }

    return (
      <div className="list-group list-special">
        <a href="#" className="list-group-item title" onClick={this.newConnectionInfo} >
          Connections
          <i className="fa fa-plus pull-right"></i>
        </a>

        <hr/>

        {bookmarks}
      </div>
    );
  }

  renderFormElements() {
    var b = this.getActiveBookmark();
    var formData = _.clone(b);

    function stringStartsWith (string, prefix) {
      return string.slice(0, prefix.length) == prefix;
    }

    if (stringStartsWith(formData["database"], initName)) {
        formData["database"] = "";
    }

    var error = "";
    var errMsg = this.state.connectionErrorMessage;
    if (errMsg !== "") {
      error = this.renderError(errMsg);
    }

    return (
      <div>
        <div className="col-md-8">
          <div className="form-group">
            <label className="control-label" htmlFor="db_hostname">Hostname</label>
            <input
              type="text"
              id="db_hostname"
              className="form-control input-sm"
              value = {formData["host"]}
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
              value = {formData["port"]}
              onChange={this.handleFormChanged.bind(this, 'port')} placeholder="5432"/>
          </div>
        </div>

        <div className="col-md-12">
          <div className="form-group">
            <label className="control-label" htmlFor="db_database">Database</label>
            <input
              type="text"
              id="db_database"
              className="form-control input-sm"
              value = {formData["database"]}
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
              value = {formData["user"]}
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
              value = {formData["password"]}
              onChange={this.handleFormChanged.bind(this, 'password')} />
          </div>
        </div>

        {error}

        <div className="col-md-12">
          <div className="form-group">
            <button disabled={this.state.isConnecting} onClick={this.handleConnect} className="btn btn-block btn-primary small">Connect</button>
            <button onClick={this.handleCancel} type="reset" id="close_connection_window" className="btn btn-block btn-default small">Cancel</button>
          </div>
        </div>
      </div>
    );
  }

  renderForm() {
    if (this.state.activeBookmark > -1) {
      var formElements = this.renderFormElements()
    } else {
      var imageStyle = {
        width: "30%",
        height: "30%"
      };

      var formElements = (
        <div className="col-md-12 text-center">
            <img class="img-responsive center-block small" src="/s/img/icon.png" alt="" style={imageStyle}/>​

            <h5>Please add a connection</h5>
        </div>
      );
    }


    return (
      <form role="form">
        {formElements}
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
        <hr/>
        <div className="connection-window-footer"><i className="fa fa-lock fa1"></i>Database crendentials are securely stored locally on your computer</div>
      </div>
    );
  }

  renderConnectionPage() {
    return (
      <div className="connection-settings">

        {this.renderConnectionWindow()}
        <hr/>

      </div>
    );
  }

  render() {
    return (
      <div id="connection_window">
          <h1>Database Workbench</h1>
          {this.renderConnectionPage()}
      </div>
    );
  }
}

module.exports = ConnectionWindow;
