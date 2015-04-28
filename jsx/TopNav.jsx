/* jshint -W097,-W117 */
'use strict';

var ViewContent = 0;
var ViewStructure = 1;
var ViewIndexes = 2;
var ViewSQLQuery = 3;
var ViewHistory = 4;
var ViewActivity = 5;
var ViewConnection = 6;

var allViews = [0, 1, 2, 3, 4, 5, 6];

var viewNames = [
  "Content",
  "Structure",
  "Indexes",
  "SQL Query",
  "History",
  "Activity",
  "Connection"
];

var TopNav = React.createClass({
  getInitialState: function() {
    return {
      view: ViewSQLQuery
    };
  },

  switchToView: function(view) {
    this.setState({
      view: view,
    });
  },

  switchToContent: function() { this.switchToView(ViewContent); },
  switchToStructure: function() { this.switchToView(ViewStructure); },
  switchToIndexes: function() { this.switchToView(ViewIndexes); },
  switchToSQLQuery: function() { this.switchToView(ViewSQLQuery); },
  switchToHistory: function() { this.switchToView(ViewHistory); },
  switchToActivity: function() { this.switchToView(ViewActivity); },
  switchToConnection: function() { this.switchToView(ViewConnection); },

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

    var children = allViews.map(function(view) {
      var cls;
      if (self.state.view == view) {
        cls = "selected";
      }
      var handler = handlers[view];
      var txt = viewNames[view];
      return <li onClick={handler} key={view} className={cls}>{txt}</li>;
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

module.exports = TopNav;
