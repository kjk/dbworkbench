import React from 'react';
import ReactDOM from 'react-dom';
import Actions from './Actions.jsx';
import * as action from './action.js';
import * as store from './store.js';

// TODO: the ace editor inside div id="custom-query" is not resized
// when we move vert drag-bar. I thought it's related to my recent
// changes but the same behavior is in 0.2.3.
// forceRender changes how we work between re-rendering the component
// or updating the style. Neither works so there's something missing
// about setting up ace editor component.
const forceRerender = true;

export default class Input extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.exportToCSV = this.exportToCSV.bind(this);
    this.handleExplain = this.handleExplain.bind(this);
    this.handleRun = this.handleRun.bind(this);

    const dy = store.getQueryEditDy();
    if (forceRerender) {
      this.state = {
        queryEditDy: dy
      };
    } else {
      this.queryEditDy = dy;
    }
  }

  componentWillMount() {
    store.onQueryEditDy((dy) => {
      if (forceRerender) {
        this.setState({
          queryEditDy: dy
        });
      } else {
        this.queryEditDy = dy;
        let el = this.editor;
        el.style.height = this.editorDy();

        el = ReactDOM.findDOMNode(this);
        el.style.height = this.inputDy();
      }
    }, this);
  }

  componentDidMount() {
    this.initEditor();
  }

  componentWillUnmount() {
    store.offAllForOwner(this);
  }

  handleRun() {
    const query = this.editor.getValue().trim();
    console.log('handleRun', query);
    if (query.length > 0) {
      action.executeQuery(query);
    }
  }

  handleExplain() {
    const query = this.editor.getValue().trim();
    console.log('handleExplain', query);
    if (query.length > 0) {
      action.explainQuery(query);
    }
  }

  exportToCSV(e) {
    e.preventDefault();
    console.log('downloadCsv');

    let query = this.editor.getValue().trim();

    if (query.length === 0) {
      return;
    }

    // Replace line breaks with spaces and properly encode query
    query = window.encodeURI(query.replace(/\n/g, ' '));

    const url = window.location.protocol + '//' + window.location.host + '/api/query?format=csv&query=' + query;
    const win = window.open(url, '_blank');
    win.focus();
  }

  initEditor() {
    this.editor = ace.edit(this.editorNode);
    this.editor.getSession().setMode('ace/mode/pgsql');
    this.editor.getSession().setTabSize(2);
    this.editor.getSession().setUseSoftTabs(true);
    this.editor.commands.addCommands([
      {
        name: 'run_query',
        bindKey: {
          win: 'Ctrl-Enter',
          mac: 'Command-Enter'
        },
        exec: (editor) => this.handleRun()
      },
      {
        name: 'explain_query',
        bindKey: {
          win: 'Ctrl-E',
          mac: 'Command-E'
        },
        exec: (editor) => this.handleExplain()
      }
    ]);
    this.editor.focus();
  }

  inputDy() {
    const dy = forceRerender ? this.state.queryEditDy : this.queryEditDy;
    return dy + 'px';
  }

  editorDy() {
    const dy = forceRerender ? this.state.queryEditDy : this.queryEditDy;
    return (dy - 50) + 'px';
  }

  render() {
    // TODO: re-add csv support
    //   <input type="button" onClick={this.exportToCSV} id="csv"
    // value="Download CSV" className="btn btn-sm btn-default" />

    //console.log("Input.render");
    let style = {};
    let editorStyle = {};

    const dy = forceRerender ? this.state.queryEditDy : this.queryEditDy;
    if (dy != 0) {
      style = {
        height: this.inputDy()
      };
      editorStyle = {
        height: this.editorDy()
      };
    }

    const nEdited = Object.keys(this.props.editedCells).length;
    const showActions = (nEdited == 0);

    return (
      <div id="input" style={ style }>
        <div className="wrapper">
          <div id="custom-query" ref={ c => this.editorNode = c } style={ editorStyle } />
          { showActions ?
            <Actions supportsExplain={ this.props.supportsExplain } onRun={ this.handleRun } onExplain={ this.handleExplain } />
            : null }
        </div>
      </div>
      );
  }
}
