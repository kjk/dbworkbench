import React from 'react';
import ReactDOM from 'react-dom';
import SpinnerCircle from './SpinnerCircle.jsx';
import DragBarHoriz from './DragBarHoriz.jsx';
import * as action from './action.js';
import * as store from './store.js';

export default class Input extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.exportToCSV = this.exportToCSV.bind(this);
    this.runExplain = this.runExplain.bind(this);
    this.runQuery = this.runQuery.bind(this);

    this.queryEditDy = store.getQueryEditDy();
  }

  componentWillMount() {
    store.onQueryEditDy((dy) => {
      this.queryEditDy = dy;

      let el = ReactDOM.findDOMNode(this.refs.editor);
      el.style.height = this.editorDy();
      //console.log('this.refs.editor.style.height: ', el.style.height);

      el = ReactDOM.findDOMNode(this);
      el.style.height = this.inputDy();
    //console.log('Input.style.height: ', el.style.height);
    }, this);
  }

  componentWillUnmount() {
    store.offAllForOwner(this);
  }

  runQuery(e) {
    if (e) {
      e.preventDefault();
    }
    const query = this.editor.getValue().trim();
    console.log('runQuery', query);
    if (query.length > 0) {
      action.executeQuery(query);
    }
  }

  runExplain(e) {
    if (e) {
      e.preventDefault();
    }
    const query = this.editor.getValue().trim();
    console.log('runExplain', query);
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
    const editorNode = ReactDOM.findDOMNode(this.refs.editor);
    this.editor = ace.edit(editorNode);
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
        exec: (editor) => this.runQuery()
      },
      {
        name: 'explain_query',
        bindKey: {
          win: 'Ctrl-E',
          mac: 'Command-E'
        },
        exec: (editor) => this.runExplain()
      }
    ]);
    this.editor.focus();
  }

  componentDidMount() {
    this.initEditor();
  }

  inputDy() {
    return this.queryEditDy + 'px';
  }

  editorDy() {
    return (this.queryEditDy - 50) + 'px';
  }

  render() {
    // TODO: re-add csv support
    //   <input type="button" onClick={this.exportToCSV} id="csv"
    // value="Download CSV" className="btn btn-sm btn-default" />

    //console.log("Input.render");
    let style = {};
    let editorStyle = {};

    if (this.queryEditDy != 0) {
      style = {
        height: this.inputDy()
      };
      editorStyle = {
        height: this.editorDy()
      };
    }

    const nEditedRows = Object.keys(this.props.editedCells).length;
    const actionsCls = nEditedRows != 0 ? 'hidden' : 'actions';
    let explainCls = 'btn btn-sm btn-default';
    if (!this.props.supportsExplain) {
      explainCls += ' hidden';
    }

    return (
      <div id="input" style={ style }>
        <div className="wrapper">
          <div id="custom-query" ref="editor" style={ editorStyle } />
          <DragBarHoriz initialY={ store.getQueryEditDy() }
            min={ 60 }
            max={ 400 }
            onPosChanged={ (dy) => store.setQueryEditDy(dy) } />
          <div className={ actionsCls }>
            <input type="button"
              onClick={ this.runQuery }
              id="run"
              value="Run Query"
              className="btn btn-sm btn-primary" />
            <input type="button"
              onClick={ this.runExplain }
              id="explain"
              value="Explain Query"
              className={ explainCls } />
            <SpinnerCircle style={ {  display: 'inline-block',  top: '4px'} } />
          </div>
        </div>
      </div>
      );
  }
}
