/* jshint -W097,-W117 */
'use strict';

var TopNav = React.createClass({

  renderSignIn: function() {
    return <a href="/logingoogle?redir=/">Sign in / Sign up</a>;
  },

  renderSignOut: function() {
    return <a href="/logout?redir=/">Sign out</a>;
  },

  render: function() {
    console.log("isLoggedIn:", this.props.isLoggedIn);
    var signInOrOut = this.props.isLoggedIn ? this.renderSignOut() : this.renderSignIn();
    return (
      <header className="navbar navbar-default navbar-fixed-top">
        <div className="container">
          <div className="navbar-header">
            <a href="/" className="navbar-brand">Database Workbench</a>
          </div>
          <nav className="collapse navbar-collapse bs-navbar-collapse">
            <ul className="nav navbar-nav">
            </ul>
            <ul className="nav navbar-nav navbar-right">
              <li>
                {signInOrOut}
              </li>
            </ul>
          </nav>
        </div>
      </header>
    );
  },
});

module.exports = TopNav;
