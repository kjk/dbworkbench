/* jshint -W097,-W117 */
'use strict';

var React = require('react');

var ReactDOM = require('react-dom');

var action = require('./action.js');

class Input extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.exportToCSV = this.exportToCSV.bind(this);
    this.runExplain = this.runExplain.bind(this);
    this.runQuery = this.runQuery.bind(this);
  }

  runQuery(e) {
    if (e) {
      e.preventDefault();
    }
    var query = $.trim(this.editor.getValue());
    console.log("runQuery", query);
    if (query.length > 0) {
      action.executeQuery(query);
    }
  }

  runExplain(e) {
    if (e) {
      e.preventDefault();
    }
    var query = $.trim(this.editor.getValue());
    console.log("runExplain", query);
    if (query.length > 0) {
      action.explainQuery(query);
    }
  }

  exportToCSV(e) {
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
  }

  initEditor() {
    var editorNode = ReactDOM.findDOMNode(this.refs.editor);
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
  }

  componentDidMount() {
    this.initEditor();
  }

  render() {
    // TODO: add csv support
    //   <input type="button" onClick={this.exportToCSV} id="csv"
    // value="Download CSV" className="btn btn-sm btn-default" />

    return (
      <div id="input">
        <div className="wrapper">
          <div id="custom-query" ref="editor"></div>
          <div className="actions">
            <input type="button" onClick={this.runQuery} id="run"
              value="Run Query" className="btn btn-sm btn-primary" />
            <input type="button" onClick={this.runExplain} id="explain"
              value="Explain Query" className="btn btn-sm btn-default" />
            <div id="query_progress">Please wait, query is executing...</div>
            <div id="input-row-number">{this.props.tooLong}</div>
          </div>
        </div>
      </div>
    );
  }
}

module.exports = Input;
