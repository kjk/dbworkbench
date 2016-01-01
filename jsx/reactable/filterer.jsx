import React from 'react';
import ReactDOM from 'react-dom';

export class Filterer extends React.Component {
  onChange() {
    this.props.onFilter(ReactDOM.findDOMNode(this).value);
  }

  render() {
    return (
      <input type="text"
        className="reactable-filter-input"
        placeholder={ this.props.placeholder }
        value={ this.props.value }
        onKeyUp={ this.onChange.bind(this) }
        onChange={ this.onChange.bind(this) } />
      );
  }
}
