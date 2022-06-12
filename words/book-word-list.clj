(defn wrap-to-json [list-name words-list]
  (json/generate-string {:metadata {:name list-name
                                    :size (count words-list)
                                    :packagedAt (java.util.Date/from (java.time.Instant/now))
                                    :version 1}

                         :words words-list} {:pretty true}))

(defn run [[book-uri list-name output-file]]
  (->> (slurp book-uri)
       (re-seq #"[\w|’|']+")
       (map clojure.string/lower-case)
       (filter #(and (> 7 (count %)) (< 1 (count %))))
       (frequencies)
       (sort-by val)
       (take-last 500)
       (map key)
       (shuffle)
       (map #(clojure.string/replace % #"’" "'"))
       (map #(clojure.string/replace % #"_" ""))
       (wrap-to-json list-name)
       (spit output-file)))

(run *command-line-args*)

(println "Done!")
