<?php

class Database
{
    private array $options;
    private PDO $pdo;

    function __construct(){

        $this->options = [
            \PDO::ATTR_ERRMODE            => \PDO::ERRMODE_EXCEPTION, //turn on errors in the form of exceptions
            \PDO::ATTR_DEFAULT_FETCH_MODE => \PDO::FETCH_ASSOC, // turn off emulation mode for "real" prepared statements
            \PDO::ATTR_EMULATE_PREPARES   => false //turn off emulation mode for "real" prepared statements
        ];
        $db = "host.docker.internal"
        $dbName = "topology-tool"

        $dsn = "mysql:host=$db;dbname="$dbName;
        $user = "www-data";
        $pwd = "password123";
        try {
            $this->pdo = new \PDO($dsn, $user, $pwd, $this->options);
        }catch (\Exception $e){
            error_log("No Connection to $dsn, with this error:\n $e.");
        }

        $this->allOccupiedDates = $this->getOccupiedDates();

    }

    //--------------------------------------------------------
    // interface function to access the database
    //--------------------------------------------------------

	// table = "tr_1", cols_string = "var1, var2", vars = [$var1,$var2]
    public function insertInto($table,$cols_string,$vars_array){
        $qm = "";
        for($i=0;$i<count($vars_array);$i++){$qm .= "?,";}
        $qm = substr($qm, 0, -1);
        $query = "INSERT INTO $table ($cols_string)VALUES($qm)";
        try{
            $statement = $this->pdo->prepare($query);
            $statement->execute($vars_array);
            return $statement->rowCount();
        }catch (\Exception $e){
            error_log($e->getMessage());
//    error_log(json_encode($e));
            return 0;
        }
    }

    protected function doesExist($table,$column,$search_var){
        $query = "SELECT * FROM $table WHERE $column = :var";
        try{
            $statement = $this->pdo->prepare($query);
            $statement->bindValue(":var",$search_var);
            $statement->execute();
            $result = $statement->fetch();
            //error_log(json_encode($result));
            $result !== false ? $response = 1 : $response = 0; // 1=exist, 0=does_not_exist
            //error_log($response);
        }catch (\Exception $e){
            error_log(json_encode($e));
        }
        return $response;
    }

    public function getId($table,$column,$search_var){
        $query = "SELECT id FROM $table WHERE $column = :var";
        try{
            $statement = $this->pdo->prepare($query);
            $statement->bindValue(":var",$search_var);
            $statement->execute();
            $result = $statement->fetch();
            //error_log(json_encode($result));
        }catch (\Exception $e){
            error_log(json_encode($e));
        }
        return $result["id"];
    }

    public function getRowById($table,$search_id){// fetch row by a certain id
        $query = "SELECT * FROM $table WHERE id = :id";
        $num_rows = self::countRows($table);
        if ($num_rows !== 0){
            try{
                $statement = $this->pdo->prepare($query);
                //error_log(json_encode($statement));
                $statement->bindValue(":id",$search_id);
                //error_log(json_encode($statement));
                $statement->execute();
                //error_log(json_encode($statement));
                $result = $statement->fetch();
            }catch (\Exception $e){
                error_log(json_encode($e));
            }
            return $result;
        }else{
            return "table_empty";
        }
//    error_log(json_encode($result));

    }

    public function countRows($table){// count rows of a table
        $query = "SELECT COUNT(*) FROM $table";
        try{
            $statement = $this->pdo->prepare($query);
            $statement->execute();
            $result = $statement->fetch();
        }catch (\Exception $e){
            error_log(json_encode($e));
        }
        //error_log(json_encode($result));
        return $result["COUNT(*)"];
    }

    public function countTables(){
        $table = function ($table){return "[$table]";};
        $query = "SHOW TABLES";
        try{
            $statement = $this->pdo->prepare($query);
            $statement->execute();
            $result = $statement->fetchAll(PDO::FETCH_FUNC,$table);
        }catch (\Exception $e){
            error_log(json_encode($e));
        }
        //error_log(json_encode($result));
        return count($result);
    }

    public function updateRow($table,$columns_array,$vars_array,$where_var,$var){
        $columns = "";
        $num_rows = $this->countRows($table);
        if($num_rows !== 0){
            foreach($columns_array as $col){$columns .= "$col=?,";}
            $columns = substr($columns, 0, -1);
            $query = "UPDATE $table SET $columns WHERE $where_var= '$var'";
            try{
                $stmt= $this->pdo->prepare($query);
                $stmt->execute($vars_array);
            }catch (\Exception $e){
                error_log(json_encode($e));
            }
        }
    }

    protected function deleteRow($table,$id){
        $query = "DELETE FROM $table WHERE id=:id";
        try{
            $statement = $this->pdo->prepare($query);
            $statement->bindValue(":id",$id);
            $statement->execute();
        }catch (\Exception $e){
            error_log(json_encode($e));
        }
    }

    public function whereColEqualsVar($table,$col_name,$search_var){// get row where col equals var
        $query = "SELECT * FROM $table WHERE $col_name=:search_var";
        try{
            $statement = $this->pdo->prepare($query);
            $statement->bindValue(":search_var",$search_var);
            $statement->execute();
            $result = $statement->fetchAll();
        }catch (\Exception $e){
            error_log(json_encode($e));
        }
        return $result;
    }

    public function setAutoincrement($table){
        $rows_num = self::countRows($table);
        $x = $rows_num+1;
        $query = "ALTER TABLE $table AUTO_INCREMENT=$x";
        try{
            $statement = $this->pdo->prepare($query);
            $statement->execute();
        }catch (\Exception $e){
            error_log(json_encode($e));
        }
    }

    function getAllRowsFromTable($table){
        $query = "SELECT * FROM $table";
        try{
            $statement = $this->pdo->prepare($query);
            $statement->execute();
            $result = $statement->fetchAll();
        }catch (\Exception $e){
            error_log(json_encode($e));
        }
        return $result;
    }

    public function getOccupiedDates(){ //  get all reservations!
        /**
         * The array Structure looks like this:
         * [[room1],[room2],[room3]]
         */

        $amountOfRoomTables = $this->countTables()-2; // minus 2 (USER and ROOM_NAMES)
        $allDates = [];

        //iterating through tables
        for ($r = 1; $r <= $amountOfRoomTables; $r++){
            $rowsInTable = $this->countRows("tr_$r");
            $roomDates = [];

            //iterating through all rows in a tables
            for($l = 1; $l <= $rowsInTable; $l++){
                //BEWARE: if the row is not existent, errors will be thrown here!
                $row = $this->getRowById("tr_$r",$l);
                // only add occupied(2) & pending(1)
                if($row["res_status"] == 2 || $row["res_status"] == 1){
                    $roomDates[] = $row;
                }
            }
            //adding all dates from one room to the final array
            $allDates[] = $roomDates;
        }
        //does give back the raw reservations without any user information
        return $allDates;
    }

    public function countRooms(){
        // minus (user-table) & (room-names-table)
        return $this->countTables() - 2;
    }

}