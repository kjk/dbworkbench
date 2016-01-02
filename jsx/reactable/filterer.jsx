import React from 'react';
import ReactDOM from 'react-dom';
import * as action from './../action.js';
import { debounce } from './../util.js';

const notifyFilterChnaged = debounce((s) => {
  console.log('notifyFilterChanged: ', s);
  action.filterChanged(s);
}, 250);

export class Filterer extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleChange = this.handleChange.bind(this);

    this.state = {
      value: this.props.defaultValue || ''
    };
  }

  componentWillMount() {
    action.onClearFilter(() => {
      this.setState({
        value: ''
      });
    }, this);
  }

  componentWillUnmount() {
    action.offAllForOwner(this);
  }

  handleChange() {
    const s = ReactDOM.findDOMNode(this).value;
    this.setState({
      value: s
    });
    notifyFilterChnaged(s);
  }

  render() {
    return (
      <input type="text"
        className="reactable-filter-input"
        placeholder={ this.props.placeholder }
        value={ this.state.value }
        onKeyUp={ this.handleChange }
        onChange={ this.handleChange } />
      );
  }
}
