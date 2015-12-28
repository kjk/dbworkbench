import React from 'react';

function pageHref(num) {
    return `#page-${num + 1}`
}

export class Paginator extends React.Component {
    constructor(props, context) {
        super(props, context);
    }

    handlePrevious(e) {
        e.preventDefault()
        this.props.onPageChange(this.props.currentPage - 1)
    }

    handleNext(e) {
        e.preventDefault()
        this.props.onPageChange(this.props.currentPage + 1);
    }

    handlePageButton(page, e) {
        e.preventDefault();
        this.props.onPageChange(page);
    }

    handleSetCurrentPage(selectedPageNumber, e) {
        if (this.props.numPages > selectedPageNumber || selectedPageNumber < 1) {
            console.log("Not possible ERRRR")
        } else {
            this.props.onPageChange(number);
        }
    }

    renderPrevious() {
        if(this.props.currentPage <= 0) {
            return (<a className='reactable-previous-page disabled'><i className="fa fa-chevron-left"></i></a>);
        }

        return (
            <a  className='reactable-previous-page'
                onClick={this.handlePrevious.bind(this)}
                href={pageHref(this.props.currentPage - 1)}>
                <i className="fa fa-chevron-left"></i>
           </a>
        );
    }

    renderNext() {
        if(this.props.currentPage >= this.props.numPages - 1) {
            return (<a className='reactable-next-page disabled'><i className="fa fa-chevron-right"></i></a>);
        }

        return (
            <a  className='reactable-next-page'
                onClick={this.handleNext.bind(this)}
                href={pageHref(this.props.currentPage + 1)}>
                <i className="fa fa-chevron-right"></i>
           </a>
        );
    }

    renderPageButton() {
        return (<span className='reactable-page'>{this.props.currentPage + 1} / {this.props.numPages} </span>);
    }

    render() {
        if (typeof this.props.numPages === 'undefined') {
            throw new TypeError('Must pass a non-zero numPages argument to Paginator');
        }

        if (typeof this.props.currentPage === 'undefined') {
            throw new TypeError('Must pass a currentPage argument to Paginator');
        }
        var st = {
          color: 'white'
        };
        
        // let pageButtons = [];
        // let pageButtonLimit = this.props.pageButtonLimit;
        // let currentPage = this.props.currentPage;
        // let numPages = this.props.numPages;
        // let lowerHalf = Math.round( pageButtonLimit / 2 );
        // let upperHalf = (pageButtonLimit - lowerHalf);

        // for (let i = 0; i < this.props.numPages; i++) {
        //     let showPageButton = false;
        //     let pageNum = i;
        //     let className = "reactable-page-button";
        //     if (currentPage === i) {
        //         className += " reactable-current-page";
        //     }
        //     pageButtons.push( this.renderPageButton(className, pageNum));
        // }

        // if(currentPage - pageButtonLimit + lowerHalf > 0) {
        //     if(currentPage > numPages - lowerHalf) {
        //         pageButtons.splice(0, numPages - pageButtonLimit)
        //     } else {
        //         pageButtons.splice(0, currentPage - pageButtonLimit + lowerHalf);
        //     }
        // }

        // if((numPages - currentPage) > upperHalf) {
        //     pageButtons.splice(pageButtonLimit, pageButtons.length - pageButtonLimit);
        // }

         // <tbody className="reactable-pagination">
         //        <tr>
         //            <td colSpan={this.props.colSpan}>
         //                <div className="reactable-pagination-button-container">
         //                {this.renderPrevious()}
         //                {this.renderPageButton()}
         //                {this.renderNext()}
         //                </div>
         //            </td>
         //        </tr>
         //    </tbody>

        return (
            <div className="reactable-pagination">
                <div className="reactable-pagination-button-container">
                    {this.renderPrevious()}
                    {this.renderPageButton()}
                    {this.renderNext()}
                    <span style={st}>{this.props.totalRowCount} rows</span>
                </div>
            </div>
        );
    }
};

