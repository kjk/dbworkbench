/* jshint -W097,-W117 */
'use strict';

var api = require('./api.js');

var ConnectionWindow = React.createClass({

  getInitialState: function() {
    return {
      connectionErrorMessage: "",
      isConnecting: false,

      standardHost: "",
      standardUsername: "",
      standardPassword: "",
      standardDatabase: "",
      standardPort: "",
      standardSSL: "require",
    };
  },

  renderBookMarks: function() {
    return (
      <div className="list-group list-special" >
        <a href="#" className="list-group-item">
          Saved Connections
          <span class="badge"><i className="fa fa-plus pull-right"></i></span>
        </a>

        <hr/>

        <a href="#" className="list-group-item active">
          <em>First Connection</em>
          <span class="badge"><i className="fa fa-times pull-right"></i></span>
        </a>
        <a href="#" className="list-group-item">
          <em>Second Connection</em>
          <span class="badge"><i className="fa fa-times pull-right"></i></span>
        </a>
        <a href="#" className="list-group-item">
          <em>Third Connection</em>
          <span class="badge"><i className="fa fa-times pull-right"></i></span>
        </a>
      </div>

    );
  },

  renderForm: function() {
    var error;
    if (this.state.connectionErrorMessage !== "") {
      error = this.renderError(this.state.connectionErrorMessage);
    }


    // <div className="col-md-4">
    //   <div className="form-group center">
    //     <div className="checkbox" name="rememberme">
    //       <label><input type="checkbox" value="" checked="checked"/>Remember Connection</label>
    //     </div>
    //   </div>
    // </div>

    return (
      <form role="form">

          <div className="col-md-8">
            <div className="form-group">
              <label className="control-label" htmlFor="db_hostname">Hostname</label>
              <input type="text" id="db_hostname" className="form-control input-sm" />
            </div>
          </div>

          <div className="col-md-4">
            <div className="form-group">
              <label className="control-label" htmlFor="db_port">Port</label>
              <input type="text" id="db_port" className="form-control input-sm" placeholder="5432"/>
            </div>
          </div>

          <div className="col-md-12">
            <div className="form-group">
              <label className="control-label" htmlFor="db_database">Database</label>
              <input type="text" id="db_database" className="form-control input-sm" />
            </div>
          </div>

          <div className="col-md-6">
            <div className="form-group">
              <label className="control-label" htmlFor="db_user">User</label>
              <input type="text" id="db_user" className="form-control input-sm" />
            </div>
          </div>

          <div className="col-md-6">
            <div className="form-group">
              <label className="control-label" htmlFor="db_pass">Password</label>
              <input type="text" id="db_pass" className="form-control input-sm" />
            </div>
          </div>
        {error}

        <div className="col-md-12">
          <div className="form-group">
            <button className="btn btn-sm btn-primary small">Connect</button>
            <button type="reset" id="close_connection_window" className="btn btn-block btn-default input-sm">Cancel</button>
          </div>
        </div>
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
