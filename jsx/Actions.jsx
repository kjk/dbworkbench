import React from 'react';
import SpinnerCircle from './SpinnerCircle.jsx';
import { Filterer } from './reactable/filterer.jsx';

export class Actions extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleRun = this.handleRun.bind(this);
    this.handleExplain = this.handleExplain.bind(this);
  }

  handleRun(e) {
    e.preventDefault();
    this.props.onRun();
  }

  handleExplain(e) {
    e.preventDefault();
    this.props.onExplain();
  }

  render() {
    return (
      <div className="actions">
        <input type="button"
          onClick={ this.handleRun }
          id="run"
          value="Run Query"
          className="btn btn-sm btn-primary" />
        { this.props.supportsExplain ?
          <input type="button"
            onClick={ this.handleExplain }
            id="explain"
            value="Explain Query"
            className='btn btn-sm btn-default' />
          : null }
        <SpinnerCircle style={ {  display: 'inline-block',  top: '4px'} } />
        <Filterer placeholder="Filter Results" defaultValue="" />
      </div>
      );
  }
}

Actions.propTypes = {
  onRun: React.PropTypes.function,
  onExplain: React.PropTypes.function,
  supportsExplain: React.PropTypes.boolean
};
