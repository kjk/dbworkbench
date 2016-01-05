import React from 'react';
import { filterInternalProps } from './utils.jsx';
import * as sort from './sort.jsx';

export class Thead extends React.Component {

  handleClickTh(columnInfo) {
    this.props.onSort(columnInfo);
  }

  renderThs() {
    const ths = this.props.columns.map((columnInfo, idx) => {
      let cls = 'reactable-th';
      if (columnInfo.isSortable) {
        cls += ' reactable-header-sortable';
      }
      if (columnInfo.sortOrder == sort.Up) {
        cls += ' reactable-header-sort-asc';
      } else if (columnInfo.sortOrder == sort.Down) {
        cls += ' reactable-header-sort-desc';
      }

      return (
        <th className={ cls }
          key={ idx }
          onClick={ this.handleClickTh.bind(this, columnInfo) }
          role="button"
          tabIndex="0">
          { columnInfo.name }
        </th>);
    });
    return ths;
  }

  render() {
    return (
      <thead>
        <tr className="reactable-column-header">
          { this.renderThs() }
        </tr>
      </thead>
      );
  }
}
