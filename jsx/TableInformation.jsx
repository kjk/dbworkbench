import React from 'react';
import filesize from 'filesize';

function isNullOrEmptyObject(object) {
  if (!object) {
    return true;
  }
  let name;
  for (name in object) {}
  return name === undefined;
}
;

export default class TableInformation extends React.Component {
  renderTableInfo(info) {
    const dataSize = parseInt(info.data_size);
    const dataSizePretty = filesize(dataSize);
    const indexSize = parseInt(info.index_size);
    const indexSizePretty = filesize(indexSize);
    const totalSize = dataSize + indexSize;
    const totalSizePretty = filesize(totalSize);
    const rowCount = parseInt(info.rows_count);

    // TODO: better done as a class,maybe on parent element
    const style = {
      backgroundColor: 'white',
    };

    return (
      <ul style={ style }>
        <li>
          <span className="table-info-light">Size:</span><span>{ totalSizePretty }</span>
        </li>
        <li>
          <span className="table-info-light">Data size:</span><span>{ dataSizePretty }</span>
        </li>
        <li>
          <span className="table-info-light">Index size:</span><span>{ indexSizePretty }</span>
        </li>
        <li>
          <span className="table-info-light">Estimated rows:</span><span>{ rowCount }</span>
        </li>
      </ul>
      );
  }

  renderTableInfoContainer() {
    const info = this.props.tableInfo;
    if (isNullOrEmptyObject(info)) {
      return;
    }

    const tableInfo = this.renderTableInfo(info);
    return (
      <div className="wrap">
        <div className="title">
          <i className="fa fa-info"></i>
          <span className="current-table-information">Table Information</span>
        </div>
        { tableInfo }
      </div>
      );
  }

  render() {
    return (
      <div className="table-information">
        { this.renderTableInfoContainer() }
      </div>
      );
  }
}

