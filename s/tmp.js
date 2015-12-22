
/* data is in format:
{
  colums: [],
  rows: [
    [], []...
  ]
}
*/

function div(s, id, cls) {
  var attr = "";

  if (id) {
    attr += ' id="' + id + '"';
  }

  if (cls) {
    attr += ' class="' + cls + '"';
  }

  return '<div' + attr + '>' + s + '</div>';
}

function td(s) {
  return '<td>' + s + '</td>\n';
}

function gen_table_html(data) {
  var s = "";
  var cols = data.columns;
  var nCols = cols.length;
  var rows = data.rows;
  var nRows = rows.length; 

  for (var i = 0; i < nCols; i++) {
    s += '<div class="results-col">';
    var colName = cols[i];
    s += div(colName);
    for (var j = 0; j < nRows; j++) {
      var row = rows[j];
      var rowVal = row[i];
      s += div(rowVal);
    }
    s += "</div>";
  }
  return div(s, "results-inner");
}

function gen_table_html2(data) {
  var s = "";
  var cols = data.columns;
  var nCols = cols.length;
  var rows = data.rows;
  var nRows = rows.length; 

  s += '<table id="results-table">';

  for (var i = 0; i < nRows; i++) {
    s += '<tr>\n';
    var row = rows[i];
    for (var j = 0; j < nCols; j++) {
      var cellVal = row[j];
      s += td(cellVal);
    }

    s += '</tr>\n';
  }
  s += '</table>';
  return s;
}