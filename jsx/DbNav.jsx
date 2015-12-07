/* jshint -W097,-W117 */
'use strict';

var action = require('./action.js');
var view = require('./view.js');

var EditConnectionButton = React.createClass({
    handleClick: function(event) {
        action.disconnectDatabase();
    },
    render: function() {
        return (
            <a href="#" id="edit_connection" className="btn btn-primary btn-xs" onClick={this.handleClick}>Disconnect</a>
        )
    }
});

var DbNav = React.createClass({
  render: function() {
    //console.log("DbNav.render: view: ", this.props.view);
    var currentView = this.props.view;
    var children = view.MainTabViews.map(function(viewName) {
      var handler = function() {
        action.viewSelected(viewName);
      };

      var selected = (currentView == viewName)
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

        <EditConnectionButton />
      </div>
    );
  },
});

module.exports = DbNav;
