import React from 'react';
import { Td } from './td.jsx';
import { toArray, filterPropsFrom } from './utils.jsx';

export class Tr extends React.Component {
  render() {
    var children = React.Children.toArray();

    if (
      this.props.data &&
      this.props.columns &&
      typeof this.props.columns.map === 'function'
    ) {
      if (typeof (children.concat) === 'undefined') {
        console.log(children);
      }

      children = children.concat(this.props.columns.map(function(column, i) {
        if (this.props.data.hasOwnProperty(column.key)) {
          var value = this.props.data[column.key];
          var props = {};

          if (
            typeof (value) !== 'undefined' &&
            value !== null &&
            value.__reactableMeta === true
          ) {
            props = value.props;
            value = value.value;
          }

          return <Td column={ column } key={ column.key } {...props}>
                   { value }
                 </Td>;
        } else {
          return <Td column={ column } key={ column.key } />;
        }
      }.bind(this)));
    }

    // Manually transfer props
    var props = filterPropsFrom(this.props);

    return React.DOM.tr(props, children);
  }
}

Tr.childNode = Td;
Tr.dataType = 'object';

