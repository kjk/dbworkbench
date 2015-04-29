/* jshint -W097,-W117 */
'use strict';

var action = require('./action.js');
var view = require('./view.js');

var TopNav = React.createClass({

  switchToView: function(view) {
    action.viewSelected(view);
  },

  switchToContent: function() { this.switchToView(view.Content); },
  switchToStructure: function() { this.switchToView(view.Structure); },
  switchToIndexes: function() { this.switchToView(view.Indexes); },
  switchToSQLQuery: function() { this.switchToView(view.SQLQuery); },
  switchToHistory: function() { this.switchToView(view.History); },
  switchToActivity: function() { this.switchToView(view.Activity); },
  switchToConnection: function() { this.switchToView(view.Connection); },

  componentDidMount: function() {
    this.switchToView(this.props.view);
  },

  render: function() {
    var handlers = [
      this.switchToContent,
      this.switchToStructure,
      this.switchToIndexes,
      this.switchToSQLQuery,
      this.switchToHistory,
      this.switchToActivity,
      this.switchToConnection
    ];
    var self = this;

    var children = view.AllViews.map(function(viewIdx) {
      var cls;
      if (self.props.view == viewIdx) {
        cls = "selected";
      }
      var handler = handlers[viewIdx];
      var txt = view.Names[viewIdx];
      return <li onClick={handler} key={viewIdx} className={cls}>{txt}</li>;
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
