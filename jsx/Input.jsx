/* jshint -W097,-W117 */
'use strict';

var Input = React.createClass({

  runQuery: function() {
    console.log("runQuery");
  },

  runExplain: function() {
    console.log("runExplain");
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
            <input type="button" id="run" value="Run Query" className="btn btn-sm btn-primary" />
            <input type="button" id="explain" value="Explain Query" className="btn btn-sm btn-default" />
            <input type="button" id="csv" value="Download CSV" className="btn btn-sm btn-default" />

            <div id="query_progress">Please wait, query is executing...</div>
          </div>
        </div>
      </div>
    );
  }
});

module.exports = Input;
