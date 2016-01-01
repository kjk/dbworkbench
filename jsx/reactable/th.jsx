import React from 'react';
import { filterPropsFrom } from './utils.jsx';

export class Th extends React.Component {
  render() {
    let childProps;

    return <th {...filterPropsFrom(this.props)}>
             { this.props.children }
           </th>;
  }
}
;

