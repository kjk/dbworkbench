import React from 'react';

function pageHref(num) {
  return `#page-${num + 1}`;
}

export default class ResultsPaginator extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handlePrevious = this.handlePrevious.bind(this);
    this.handleNext = this.handleNext.bind(this);
  }

  handlePrevious(e) {
    e.preventDefault();
    this.props.onPageChange(this.props.currentPage - 1);
  }

  handleNext(e) {
    e.preventDefault();
    this.props.onPageChange(this.props.currentPage + 1);
  }

  renderPrevious() {
    if (this.props.currentPage <= 0) {
      return (<a className='reactable-previous-page disabled'><i className='fa fa-chevron-left'/></a>);
    }

    return (
      <a className='reactable-previous-page' onClick={ this.handlePrevious } href={ pageHref(this.props.currentPage - 1) }><i className='fa fa-chevron-left'/></a>
      );
  }

  renderNext() {
    if (this.props.currentPage >= this.props.nPages - 1) {
      return (<a className='reactable-next-page disabled'><i className='fa fa-chevron-right'/></a>);
    }

    return (
      <a className='reactable-next-page' onClick={ this.handleNext } href={ pageHref(this.props.currentPage + 1) }><i className='fa fa-chevron-right'/></a>
      );
  }

  render() {
    const style = {
      color: 'white'
    };

    return (
      <div className='reactable-pagination'>
        <div className='reactable-pagination-button-container'>
          { this.renderPrevious() }
          <span className='reactable-page'>{ this.props.currentPage + 1 } / { this.props.nPages }</span>
          { this.renderNext() }
          <span style={ style }>{ this.props.nRows } rows</span>
        </div>
      </div>
      );
  }
}

ResultsPaginator.propTypes = {
  onPageChange: React.PropTypes.func.isRequired,
  currentPage: React.PropTypes.number.isRequired,
  nPages: React.PropTypes.number.isRequired,
  nRows: React.PropTypes.number.isRequired
};
