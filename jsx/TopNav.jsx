/* jshint -W097,-W117 */
'use strict';

var TopNav = React.createClass({

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
          </nav>
        </div>
      </header>
    );
  },
});

module.exports = TopNav;
