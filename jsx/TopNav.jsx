/* jshint -W097,-W117 */
'use strict';

var action = require('./action.js');
var view = require('./view.js');

var TopNav = React.createClass({

  render: function() {
    var currentView = this.props.view;
    var children = view.AllViews.map(function(viewIdx) {
      var txt = view.Names[viewIdx];
      var handler = function() {
        action.viewSelected(viewIdx);
      };
      var cls = (currentView == viewIdx) ? "selected" : "";
      return <li key={viewIdx} onClick={handler} className={cls}>{txt}</li>;
    });

    return (
      <div id="nav">
        <ul>
          {children}
        </ul>

        <a href="#" id="edit_connection" className="btn btn-primary btn-sm"><i className="fa fa-gear"></i> Edit Connection</a>
      </div>
    );
  },
});

module.exports.TopNav = TopNav;
