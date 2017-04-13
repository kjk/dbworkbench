import React from 'react';
import ReactDOM from 'react-dom';
import PropTypes from 'prop-types';
import Popover from 'react-popover';
import * as action from './action';
import * as store from './store';

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

  componentWillMount() {
    store.onQueryEditDy((dy) => {
      this.queryEditDy = dy;

      const top = this.topPos();
      this.btnSaveNode.style.top = top;
      this.btnDiscardNode.style.top = top;
      this.rowCountNode.style.top = top;
      this.sqlPreviewNode.style.top = top;
    }, this);
  }

  componentWillUnmount() {
    store.offAllForOwner(this);
  }

  handleToggleSQLPreview() {
    //console.log('handleSQLPreview');
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

  handleSaveChanges() {
    console.log('handleSaveChanges ');

    // TODO: must support multiple queries for multiple rows changes
    var query = this.props.generateQuery();
    action.executeQuery(query);
  }

  togglePopover() {
    this.setState({
      isOpen: !this.state.isOpen
    });
  }

  topPos() {
    return this.queryEditDy + 'px';
  }

  rememberNode(el) {
    console.log('rememberNode: ', el.id);

    switch (el.id) {
      case 'id-save':
        this.btnSaveNode = el;
        break;
      case 'id-discard':
        this.btnDiscardNode = el;
        break;
      case 'id-row-number':
        this.rowCountNode = el;
        break;
      case 'id-sql-preview':
        this.sqlPreviewNode = el;
        break;
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
    };

    return (
      <div id='query_edit_bar'>
        <button ref={ this.rememberNode }
          className='save_changes'
          id='id-save'
          onClick={ this.handleSaveChanges }
          style={ style }>
          Save Changes
        </button>
        <button ref={ this.rememberNode }
          id='id-discard'
          className='discard_changes'
          onClick={ this.props.onHandleDiscardChanges }
          style={ style }>
          Discard Changes
        </button>
        <div ref={ this.rememberNode }
          className='row_number'
          id='id-row-number'
          style={ style }>
          { this.props.numberOfRowsEdited } edited rows
        </div>
        <Popover style={ popOverStyle }
          isOpen={ this.state.isOpen }
          body={ this.state.popOverText }
          preferPlace={ "right" }
          target={ "sql_preview" }
          targetElement={ "sql_preview" }
          tipSize={ 10 }>
          <div ref={ this.rememberNode }
            className='sql_preview'
            id='id-sql-preview'
            onClick={ this.handleToggleSQLPreview }
            style={ style }>
            { !this.state.isOpen ? 'Show SQL Preview' :
              'Hide SQL Preview' }
          </div>
        </Popover>
      </div>
      );
  }
}

QueryEditBar.propTypes = {
  generateQuery: PropTypes.func.isRequired,
  onHandleDiscardChanges: PropTypes.func.isRequired,
  numberOfRowsEdited: PropTypes.number.isRequired
};
