/* jshint -W097,-W117 */
'use strict';

var TopNav = React.createClass({

  handleSignIn: function(e) {
    e.preventDefault();
    console.log("clicked sign in");
  },

  render: function() {

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
                <a href="#" onClick={this.handleSignIn}>Sign in / Sign up</a>
              </li>
            </ul>
          </nav>
        </div>
      </header>
    );
  },
});

module.exports = TopNav;
