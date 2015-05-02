/* jshint -W097,-W117 */
'use strict';

var action = require('./action.js');

var Input = React.createClass({

  runQuery: function(e) {
    if (e) {
      e.preventDefault();
    }
    var query = $.trim(this.editor.getValue());
    console.log("runQuery", query);
    if (query.length > 0) {
      action.executeQuery(query);
    }
  },

  runExplain: function(e) {
    if (e) {
      e.preventDefault();
    }
    var query = $.trim(this.editor.getValue());
    console.log("runExplain", query);
    if (query.length > 0) {
      action.explainQuery(query);
    }
  },

  exportToCSV: function(e) {
    e.preventDefault();
    console.log("downloadCsv");

    var query = $.trim(this.editor.getValue());

    if (query.length === 0) {
      return;
    }

    // Replace line breaks with spaces and properly encode query
    query = window.encodeURI(query.replace(/\n/g, " "));

    var url = window.location.protocol + "//" + window.location.host + "/api/query?format=csv&query=" + query;
    var win = window.open(url, '_blank');
    win.focus();
  },

  initEditor: function() {
    var editorNode = React.findDOMNode(this.refs.editor);
    this.editor = ace.edit(editorNode);
    this.editor.getSession().setMode("ace/mode/pgsql");
    this.editor.getSession().setTabSize(2);
    this.editor.getSession().setUseSoftTabs(true);

    var self = this;
    this.editor.commands.addCommands([
      {
        name: "run_query",
        bindKey: {
          win: "Ctrl-Enter",
          mac: "Command-Enter"
        },
        exec: function(editor) {
          self.runQuery();
        }
      },
      {
        name: "explain_query",
        bindKey: {
          win: "Ctrl-E",
          mac: "Command-E"
        },
        exec: function(editor) {
          self.runExplain();
        }
      }
    ]);
    this.editor.focus();
  },

  componentDidMount: function() {
    this.initEditor();
  },

  render: function() {
    return (
      <div id="input">
        <div className="wrapper">
          <div id="custom_query" ref="editor"></div>
          <div className="actions">
            <input type="button" onClick={this.runQuery} id="run"
              value="Run Query" className="btn btn-sm btn-primary" />
            <input type="button" onClick={this.runExplain} id="explain"
              value="Explain Query" className="btn btn-sm btn-default" />
            <input type="button" onClick={this.exportToCSV} id="csv"
              value="Download CSV" className="btn btn-sm btn-default" />
            <div id="query_progress">Please wait, query is executing...</div>
          </div>
        </div>
      </div>
    );
  }
});

module.exports = Input;
