/* jshint -W097,-W117 */
'use strict';

var api = require('./api.js');

var ConnectionScheme = 0;
var ConnectionStandard = 1;
var ConnectionSSH = 2;

var ConnectionWindow = React.createClass({

  getInitialState: function() {
    return {
      connectionType: ConnectionScheme,
      connectionErrorMessage: "",
      isConnecting: false,

      schemeURL: "",

      standardHost: "",
      standardUsername: "",
      standardPassword: "",
      standardDatabase: "",
      standardPort: "",
      standardSSL: "require",
    };
  },

  handleChangeConnectionTypeToScheme: function(e) {
    console.log("handleChangeConnectionTypeToScheme");
    e.preventDefault();
    this.setState({
      connectionType: ConnectionScheme,
      connectionErrorMessage: ""
    });
  },

  handleChangeConnectionTypeToStandard: function(e) {
    console.log("handleChangeConnectionTypeToStandard");
    e.preventDefault();
    this.setState({
      connectionType: ConnectionStandard,
      connectionErrorMessage: ""
    });
  },

  handleChangeConnectionTypeToSSH: function(e) {
    console.log("handleChangeConnectionTypeToSSH");
    e.preventDefault();
    this.setState({
      connectionType: ConnectionSSH,
      connectionErrorMessage: ""
    });
  },

  handleConnectStandardChanged: function(name, e) {
    var change = {};
    change[name] = e.target.value;
    console.log("handleConnectStandardChanged: ", change);
    this.setState(change);
  },

  handleConnectionSchemeChanged: function(e) {
    var scheme = e.target.value;
    console.log("handleConnectionSchemeChanged: ", scheme);
    this.setState({
      schemeURL: scheme
    });
  },

  handleConnect: function(e) {
    e.preventDefault();
    console.log("handleConnect");

    var url = this.getConnectionString();
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
      }
      else {
        console.log("did connect");
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

  getConnectionString: function() {
    var url;

    switch (this.state.connectionType) {
      case ConnectionScheme:
        url = this.state.schemeURL.trim();
        break;
      case ConnectionStandard:
        var host = this.state.standardHost;
        var port = this.state.standardPort;
        var user = this.state.standardUsername;
        var pass = this.state.standardPassword; // TODO: encoding?
        var db = this.state.standardDatabase;
        var ssl = this.state.standardSSL;

        if (port.length == 0) {
          port = "5432";
        }

        url = "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + db + "?sslmode=" + ssl;

        break;
      default:
        console.log("This shouldn't happen");
    }

    return url;
  },

  renderConnectionFormHeader: function() {
    var clsScheme = "btn btn-default";
    var clsStandard = "btn btn-default";
    var clsSSH = "btn btn-default";
    switch (this.state.connectionType) {
      case ConnectionScheme:
        clsScheme += ' active';
        break;
      case ConnectionStandard:
        clsStandard += ' active';
        break;
      case ConnectionSSH:
        clsSSH += ' active';
        break;
    }

    // TODO: add SSH section
    // <button onClick={this.handleChangeConnectionTypeToSSH} className={clsSSH}>SSH</button>

    return (
        <div className="text-center">
          <div className="btn-group btn-group-sm connection-group-switch">
            <button onClick={this.handleChangeConnectionTypeToScheme} className={clsScheme}>Scheme</button>
            <button onClick={this.handleChangeConnectionTypeToStandard} className={clsStandard} >Standard</button>
          </div>
        </div>
    );
  },

  renderStandardGroup: function() {
    return (
      <div className="connection-standard-group">
        <div className="form-group bookmarks">
          <label className="col-sm-3 control-label">Bookmark</label>
          <div className="col-sm-9">
            <select className="form-control" id="connection_bookmarks"></select>
          </div>
        </div>

        <div className="form-group">
          <label className="col-sm-3 control-label">Host</label>
          <div className="col-sm-9">
            <input type="text" id="pg_host" value={this.state.standardHost} onChange={this.handleConnectStandardChanged.bind(this, 'standardHost')} className="form-control" />
          </div>
        </div>

        <div className="form-group">
          <label className="col-sm-3 control-label">Username</label>
          <div className="col-sm-9">
            <input type="text" id="pg_user" value={this.state.standardUsername} onChange={this.handleConnectStandardChanged.bind(this, 'standardUsername')} className="form-control" />
          </div>
        </div>

        <div className="form-group">
          <label className="col-sm-3 control-label">Password</label>
          <div className="col-sm-9">
            <input type="password" id="pg_password" value={this.state.standardPassword} onChange={this.handleConnectStandardChanged.bind(this, 'standardPassword')} className="form-control" />
          </div>
        </div>

        <div className="form-group">
          <label className="col-sm-3 control-label">Database</label>
          <div className="col-sm-9">
            <input type="text" id="pg_db" value={this.state.standardDatabase} onChange={this.handleConnectStandardChanged.bind(this, 'standardDatabase')} className="form-control" />
          </div>
        </div>

        <div className="form-group">
          <label className="col-sm-3 control-label">Port</label>
          <div className="col-sm-9">
            <input type="text" id="pg_port" value={this.state.standardPort} onChange={this.handleConnectStandardChanged.bind(this, 'standardPort')} className="form-control" placeholder="5432" />
          </div>
        </div>

        <div className="form-group">
          <label className="col-sm-3 control-label">SSL</label>
          <div className="col-sm-9">
            <select className="form-control" id="connection_ssl" value={this.state.standardSSL} onChange={this.handleConnectStandardChanged.bind(this, 'standardSSL')} >
              <option value="disable">disable</option>
              <option value="require">require</option>
              <option value="verify-full">verify-full</option>
            </select>
          </div>
        </div>
      </div>
    );
  },

  renderSchemeGroup: function() {
    return (
      <div className="connection-scheme-group">
        <div className="form-group">
          <div className="col-sm-12">
            <label>Enter server URL scheme</label>
            <input type="text" value={this.state.schemeURL} onChange={this.handleConnectionSchemeChanged} className="form-control"/>
            <p className="help-block">URL format: postgres://user:password@host:port/db?sslmode=mode
            </p>
            <p className="help-block">Test database: postgres://localhost/world</p>
          </div>
        </div>
      </div>
    );
  },

  renderSSHGroup: function() {
    return (
      <div className="connection-ssh-group">
        <div className="form-group">
          <label className="col-sm-3 control-label">SSH Host</label>
          <div className="col-sm-9">
            <input type="text" id="ssh_host" className="form-control" />
          </div>
        </div>

        <div className="form-group">
          <label className="col-sm-3 control-label">SSH User</label>
          <div className="col-sm-9">
            <input type="text" id="ssh_user" className="form-control" />
          </div>
        </div>

        <div className="form-group">
          <label className="col-sm-3 control-label">SSH Password</label>
          <div className="col-sm-9">
            <input type="text" id="ssh_password" className="form-control" placeholder="optional" />
          </div>
        </div>

        <div className="form-group">
          <label className="col-sm-3 control-label">SSH Port</label>
          <div className="col-sm-9">
            <input type="text" id="pg_host" className="form-control" placeholder="optional" />
          </div>
        </div>
      </div>
    );
  },

  renderError: function(errorText) {
    return (
      <div className="alert alert-danger">{errorText}</div>
    );
  },

  renderConnectionWindow: function() {
    var connectionFormHeader = this.renderConnectionFormHeader();
    var group;
    switch (this.state.connectionType) {
      case ConnectionScheme:
        group = this.renderSchemeGroup();
        break;
      case ConnectionStandard:
        group = this.renderStandardGroup();
        break;
      case ConnectionSSH:
        group = this.renderSSHGroup();
        break;
      default:
        console.log("unknown connectionType: ", this.state.connectionType);
    }
    var error;
    if (this.state.connectionErrorMessage !== "") {
      error = this.renderError(this.state.connectionErrorMessage);
    }
    var connectDisabled = this.state.isConnecting;

    return (
      <div id="connection_window">
        <div className="connection-settings">
          <h1>Postgres Database Workbench</h1>
            <form role="form" className="form-horizontal" id="connection_form">
            {connectionFormHeader}
            <hr/>

            {group}
            {error}

            <div className="form-group">
              <div className="col-sm-12">
                <button disabled={connectDisabled} onClick={this.handleConnect} className="btn btn-block btn-primary">Connect</button>
                <button onClick={this.handleCancel} type="reset" id="close_connection_window" className="btn btn-block btn-default">Cancel</button>
              </div>
            </div>
          </form>
        </div>
      </div>
    );
  },

  renderPleaseSignIn: function() {
    return (
      <div id="connection_window">
        <div className="connection-settings">
          <h1>Postgres Database Workbench</h1>
          <p>Please sign in to connect to a database.</p>
        </div>
      </div>
    );
  },

  render: function() {
    return (
      <div>
        {this.renderConnectionWindow()}
      </div>
    );
  }
});

module.exports = ConnectionWindow;
