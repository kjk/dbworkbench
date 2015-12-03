/* jshint -W097,-W117 */
'use strict';

var api = require('./api.js');
var action = require('./action.js');
var _ = require('underscore');

var initName = "Unnamed Connection";

var ConnectionWindow = React.createClass({

  getInitialState: function() {
    // TODO: need a solution for this not making apicall sync
    var bookmarks = {};
    api.getBookmarks(function(data) {
      console.log("getBookmarks: ", data);
      if (data != undefined && data["error"] == null) {
        bookmarks = data;
      }
    });

    var activeBookmark = "";
    if (Object.keys(bookmarks).length !== 0) {
      activeBookmark = Object.keys(bookmarks)[0];
    }

    return {
      connectionErrorMessage: "",
      isConnecting: false,

      bookmarks: bookmarks,
      activeBookmark: activeBookmark,
    };
  },

  addBookmark: function(e) {
    var bookmarkLimit = 10;
    if (Object.keys(this.state.bookmarks).length >= bookmarkLimit) {
      action.alertBar("Max Connection Limit is " + bookmarkLimit);
      return;
    }

    var possibleNames = [];
    possibleNames.push(initName);
    for (var i = 1; i <= bookmarkLimit; i++) {
      possibleNames.push(initName + " " + String(i));
    };

    var bookmarkTitles = [];
    for (var bookmarkName in this.state.bookmarks) {
      bookmarkTitles.push(bookmarkName);
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
        activeBookmark: newName,
      });
    });
  },

  getBookmark: function(e) {
    var self = this;
    api.getBookmarks(function(data) {
      console.log("getBookmarks: ", data);

      self.setState({
        bookmarks: data,
      });
    });
  },

  deleteBookmark: function(e) {
    e.stopPropagation();
    var self = this;
    var dbName = event.target.attributes["id"].value;
    // console.log(dbName)
    // console.log(event.target.attributes["id"])
    api.removeBookmark(dbName, function(data) {
      console.log("removeBookmarks: ", data);

      if (data !== undefined && Object.keys(data).length > 0) {
        var activeBookmark = Object.keys(data)[0];
        var bookmarks = data;
      } else {
        var bookmarks = {};
        var activeBookmark = "";
      }

      self.setState({
        bookmarks: bookmarks,
        activeBookmark: activeBookmark,
      });
    });
  },

  selectBookmark: function(e) {
    console.log("selectBookmark");

    this.setState({
      activeBookmark: event.target.attributes["id"].value,
    });
  },

  handleFormChanged: function(name, e) {
    console.log("handleFormChanged: ", e.target.value);

    var change = this.state.bookmarks;
    change[this.state.activeBookmark][name] = e.target.value;

    this.setState({
      bookmarks: change,
      connectionErrorMessage: "",
    });
  },

  handleConnect: function(e) {
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

        api.addBookmark(self.state.bookmarks[self.state.activeBookmark], function(data) {
          console.log("bookmark saved: ", data);
        });

        var connId = resp.ConnectionID;
        var connStr = url;
        var databaseName = resp.CurrentDatabase;
        self.props.onDidConnect(connStr, connId, databaseName);

      }
    });
  },

  handleCancel: function(e) {
    e.preventDefault();
    console.log("handleCancel");
  },

  renderError: function(errorText) {
    return (
      <div className="alert alert-danger">{errorText}</div>
    );
  },

  renderBookMarks: function() {
    var bookmarks = [];
    for (var bookmarkName in this.state.bookmarks) {
      var databaseName = this.state.bookmarks[bookmarkName]["database"];
      var removeButton = <i id={bookmarkName} onClick={this.deleteBookmark} className="fa fa-times pull-right"></i>;

      var className = "list-group-item"
      if (bookmarkName == this.state.activeBookmark) {
        className = "list-group-item active"
      }

      bookmarks.push(
        <a id={bookmarkName} href="#" className={className} onClick={this.selectBookmark}>
          <em id={bookmarkName}>{databaseName}</em>
          {removeButton}
        </a>
      );
    }

    return (
      <div className="list-group list-special">
        <a href="#" className="list-group-item">
          Connection List
          <i id={bookmarkName} onClick={this.addBookmark} className="fa fa-plus pull-right"></i>
        </a>

        <hr/>

        {bookmarks}
      </div>

    );
  },

  renderFormElements: function() {

    var formData = {};
    if (Object.keys(this.state.bookmarks).length !== 0) {
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
  },

  renderForm: function() {
    if (this.state.activeBookmark !== "") {
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
  },

  renderConnectionWindow: function() {
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


        <h6>Database crendentials are securely stored locally on your computer</h6>

      </div>
    );
  },

  renderConnectionPage: function() {
    return (
      <div className="connection-settings">

        {this.renderConnectionWindow()}
        <hr/>

      </div>
    );
  },

  render: function() {
    return (
      <div id="connection_window">
          <h1>Postgres Database Workbench</h1>

          <hr/>

          {this.renderConnectionPage()}
      </div>
    );
  }
});

module.exports = ConnectionWindow;
