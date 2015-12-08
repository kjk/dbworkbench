/* jshint -W097,-W117 */
'use strict';

var React = require('react');

var api = require('./api.js');
var action = require('./action.js');
var _ = require('underscore');

var initName = "Unnamed Connection";

// string.startsWith is es6
// TODO: shouldn't this be taken care by es6 transform?
if (!String.prototype.startsWith) {
  String.prototype.startsWith = function(searchString, position) {
    position = position || 0;
    return this.indexOf(searchString, position) === position;
  };
}

class ConnectionWindow extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.addBookmark = this.addBookmark.bind(this);
    this.deleteBookmark = this.deleteBookmark.bind(this);
    this.handleCancel = this.handleCancel.bind(this);
    this.handleConnect = this.handleConnect.bind(this);
    this.handleFormChanged = this.handleFormChanged.bind(this);
    this.selectBookmark = this.selectBookmark.bind(this);

    this.state = {
      connectionErrorMessage: "",
      isConnecting: false,

      bookmarks: gBookmarkInfo,
      activeBookmark: (gBookmarkInfo.length !== 0) ? 0 : -1,
    };
  }

  findBookmarkByDatabaseName(databaseName) {
    var index = _.findIndex(this.state.bookmarks, function(obj) {
      return (obj["database"] == databaseName)
    });

    return index;
  }

  addBookmark(e) {
    var bookmarkLimit = 10;
    if (this.state.bookmarks.length >= bookmarkLimit) {
      action.alertBar("Max Connection Limit is " + bookmarkLimit); // TODO: this doesn't work check why
      return;
    }

    var possibleNames = [];
    possibleNames.push(initName);
    for (var i = 1; i <= bookmarkLimit; i++) {
      possibleNames.push(initName + " " + String(i));
    };

    var bookmarkTitles = [];
    var n = this.state.bookmarks.length;
    for (var i = 0; i < n; i++ ) {
      var bookmark = this.state.bookmarks[i];
      bookmarkTitles.push(bookmark["database"]);
    }

    var usedNames = _.intersection(possibleNames, bookmarkTitles);
    var notUsedNames = _.difference(possibleNames, usedNames);

    if (notUsedNames.length > 0 ) {
      var newName = notUsedNames[0];
    } else {
      console.log("Error: this should not be possible");
      return;
    }

    console.log("new bookmark name" + newName);

    var initialBookmark = { url: "",
                            host: "",
                            database: newName,
                            user: "",
                            password: "",
                            port: "" ,
                            ssl: ""
                          };

    var self = this;
    api.addBookmark(initialBookmark, function(data) {
      console.log("bookmark added: ", data);

      self.setState({
        bookmarks: data,
        activeBookmark: self.findBookmarkByDatabaseName(newName),
      });
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

    var dbName = e.target.attributes["data-custom-attribute"].value;
    api.removeBookmark(dbName, function(data) {
      console.log("deleteBookmarks removing: ", dbName, " data: ", data);

      var bookmarks = [];
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
    console.log("selectBookmark", e.target.attributes["id"].value);

    this.setState({
      activeBookmark: e.target.attributes["id"].value,
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

  handleConnect(e) {
    e.preventDefault();
    console.log("handleConnect");

    var host = this.state.bookmarks[this.state.activeBookmark]["host"];
    var port = this.state.bookmarks[this.state.activeBookmark]["port"];
    var user = this.state.bookmarks[this.state.activeBookmark]["user"];
    var pass = this.state.bookmarks[this.state.activeBookmark]["password"];
    var db = this.state.bookmarks[this.state.activeBookmark]["database"];
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
        console.log("did connect");

        api.removeBookmark(self.state.bookmarks[self.state.activeBookmark]["oldDatabase"], function(data) {
          api.addBookmark(self.state.bookmarks[self.state.activeBookmark], function(data) {
            console.log("bookmark saved: ", data);
          });
        });

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
    for (var position = 0; position < this.state.bookmarks.length; position++) {
      var bookmark = this.state.bookmarks[position];
      var databaseName = bookmark["database"];

      var removeButton = <i data-custom-attribute={databaseName} onClick={this.deleteBookmark} className="fa fa-times pull-right"></i>;

      var className = "list-group-item"
      if (position == this.state.activeBookmark) {
        className = "list-group-item active"
      }

      bookmarks.push(
        <a id={position} key={position} href="#" className={className} onClick={this.selectBookmark}>
          {databaseName}
          {removeButton}
        </a>
      );
    }

    return (
      <div className="list-group list-special">
        <a href="#" className="list-group-item title" onClick={this.addBookmark} >
          Connections
          <i className="fa fa-plus pull-right"></i>
        </a>

        <hr/>

        {bookmarks}
      </div>

    );
  }

  renderFormElements() {
    var formData = [];
    if (this.state.bookmarks.length > 0) {
      formData = _.clone(this.state.bookmarks[this.state.activeBookmark])

      if (formData["database"].startsWith(initName)) {
        formData["database"] = "";
      }
    } else {
      formData["database"] = "";
    }

    if (this.state.connectionErrorMessage !== "") {
      var error = this.renderError(this.state.connectionErrorMessage);
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
          <h1>Postgres Database Workbench</h1>
          {this.renderConnectionPage()}
      </div>
    );
  }
}

module.exports = ConnectionWindow;
