/* jshint -W097,-W117 */
'use strict';

var utils = require('./utils.js');

var App = React.createClass({
  getInitialState: function() {
    return {
      connectionId: null,
    }
  },

  renderNoConnection: function() {
    return (
      <div id="connection_window">
        <div className="connection-settings">
          <h1>Database Workbench</h1>

          <form role="form" className="form-horizontal" id="connection_form">
            <div className="text-center">
              <div className="btn-group btn-group-sm connection-group-switch">
                <button type="button" data="scheme" className="btn btn-default" id="connection_scheme">Scheme</button>
                <button type="button" data="standard" className="btn btn-default active" id="connection_standard">Standard</button>
              </div>
            </div>

            <hr/>

            <div className="connection-scheme-group">
              <div className="form-group">
                <div className="col-sm-12">
                  <label>Enter server URL scheme</label>
                  <input type="text" className="form-control" id="connection_url" name="url" />
                  <p className="help-block">URL format: postgres://user:password@host:port/db?sslmode=mode
                  </p>
                </div>
              </div>
            </div>

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
                  <input type="text" id="pg_host" className="form-control" />
                </div>
              </div>

              <div className="form-group">
                <label className="col-sm-3 control-label">Username</label>
                <div className="col-sm-9">
                  <input type="text" id="pg_user" className="form-control" />
                </div>
              </div>

              <div className="form-group">
                <label className="col-sm-3 control-label">Password</label>
                <div className="col-sm-9">
                  <input type="text" id="pg_password" className="form-control" />
                </div>
              </div>

              <div className="form-group">
                <label className="col-sm-3 control-label">Database</label>
                <div className="col-sm-9">
                  <input type="text" id="pg_db" className="form-control" />
                </div>
              </div>

              <div className="form-group">
                <label className="col-sm-3 control-label">Port</label>
                <div className="col-sm-9">
                  <input type="text" id="pg_port" className="form-control" placeholder="5432" />
                </div>
              </div>

              <div className="form-group">
                <label className="col-sm-3 control-label">SSL</label>
                <div className="col-sm-9">
                  <select className="form-control" id="connection_ssl" defaultValue="require">
                    <option value="disable">disable</option>
                    <option value="require">require</option>
                    <option value="verify-full">verify-full</option>
                  </select>
                </div>
              </div>
            </div>

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

            <div id="connection_error" className="alert alert-danger"></div>

            <div className="form-group">
              <div className="col-sm-12">
                <button type="submit" className="btn btn-block btn-primary">Connect</button>
                <button type="reset" id="close_connection_window" className="btn btn-block btn-default">Cancel</button>
              </div>
            </div>
          </form>
        </div>
      </div>
    );
  },

  render: function() {
    if (this.state.connectionId === null) {
      return this.renderNoConnection();
    } else {
      return <div>This is a start</div>;
    }
  }
});

function appStart() {
  React.render(
    <App/>,
    document.getElementById('main')
  );
}

window.appStart = appStart;
