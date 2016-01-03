import React from 'react';
import ReactDOM from 'react-dom';
import Popover from 'react-popover';
import * as action from './action.js';
import * as store from './store.js';

export default class QueryEditBar extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleSaveChanges = this.handleSaveChanges.bind(this);
    this.handleToggleSQLPreview = this.handleToggleSQLPreview.bind(this);

    // 1) Is there a way to move discard changes to here without using action?
    // 2) maybe move generateQuery from output to here?

    this.queryEditDy = store.getQueryEditDy();

    this.state = {
      isOpen: false,
      popOverText: '',
    };
  }

  setTopOnRef(ref, top) {
    const el = ReactDOM.findDOMNode(ref);
    el.style.top = top;
  }

  componentWillMount() {
    store.onQueryEditDy((dy) => {
      this.queryEditDy = dy;

      const top = this.topPos();
      this.setTopOnRef(this.refs.btnSave, top);
      this.setTopOnRef(this.refs.btnDiscard, top);
      this.setTopOnRef(this.refs.rowCount, top);
      this.setTopOnRef(this.refs.sqlPreview, top);
    }, this);
  }

  componentWillUnmount() {
    store.offAllForOwner(this);
  }

  togglePopover() {
    this.setState({
      isOpen: !this.state.isOpen
    });
  }

  handleSaveChanges() {
    console.log('handleSaveChanges ');

    // TODO: must support multiple queries for multiple rows changes
    var query = this.props.generateQuery();
    action.executeQuery(query);
  }

  topPos() {
    return this.queryEditDy + 'px';
  }

  handleToggleSQLPreview() {
    console.log('handleSQLPreview');
    if (this.state.isOpen) {
      this.setState({
        isOpen: false
      });
    } else {
      var query = this.props.generateQuery();
      query = query.split(';').join('\n');

      this.setState({
        popOverText: query,
        isOpen: true,
      });
    }
  }

  render() {
    // TODO: try setting top on query_edit_bar element
    // instead of on each child
    var style = {
      top: this.topPos()
    };

    var popOverStyle = {
      zIndex: '4',
    }

    return (
      <div id="query_edit_bar">
        <button ref="btnSave"
          className="save_changes"
          onClick={ this.handleSaveChanges.bind(this) }
          style={ style }>
          Save Changes
        </button>
        <button ref="btnDiscard"
          className="discard_changes"
          onClick={ this.props.onHandleDiscardChanges }
          style={ style }>
          Discard Changes
        </button>
        <div ref="rowCount" className="row_number" style={ style }>
          { this.props.numberOfRowsEdited } edited rows
        </div>
        <Popover
          style={popOverStyle}
          isOpen={ this.state.isOpen }
          body={ this.state.popOverText }
          preferPlace={ "right" }
          target={ "sql_preview" }
          targetElement={ "sql_preview" }
          tipSize={ 10 }>
          <div ref="sqlPreview"
            className="sql_preview"
            onClick={ this.handleToggleSQLPreview.bind(this) }
            style={ style }>
            { !this.state.isOpen ? 'Show SQL Preview' :
              'Hide SQL Preview' }
          </div>
        </Popover>
      </div>
      );
  }
}
