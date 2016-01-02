import React from 'react';
import ReactDOM from 'react-dom';
import * as action from './../action.js'

export class Filterer extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleChange = this.handleChange.bind(this);
  }

  handleChange() {
    const s = ReactDOM.findDOMNode(this).value;
    action.filterChanged(s);
  }

  render() {
    return (
      <input type="text"
        className="reactable-filter-input"
        placeholder={ this.props.placeholder }
        value={ this.props.value }
        onKeyUp={ this.handleChange }
        onChange={ this.handleChange } />
      );
  }
}
