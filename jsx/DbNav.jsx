/* jshint -W097,-W117 */
'use strict';

var React = require('react');

var action = require('./action.js');
var view = require('./view.js');

class DbNav extends React.Component {
  render() {
    //console.log("DbNav.render: view: ", this.props.view);
    var currentView = this.props.view;
    var children = view.MainTabViews.map(function(viewName) {
      var handler = function() {
        action.viewSelected(viewName);
      };

      var selected = (currentView == viewName);
      if (selected) {
        return <li key={viewName} onClick={handler} className="selected"><u>{viewName}</u></li>;
      } else {
        return <li key={viewName} onClick={handler}>{viewName}</li>;
      }
    });

    return (
      <div id="nav">
        <ul>
          {children}
        </ul>
      </div>
    );
  }
}

module.exports = DbNav;
